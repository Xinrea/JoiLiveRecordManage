package api

import (
	"errors"
	"fmt"
	"joirecord/internal/db"
	"joirecord/internal/logger"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var log = logger.Log.WithField("module", "apiServer")

type Server struct {
	s     *gin.Engine
	c     *bos.Client
	cache *RecordCache
}

func New(bosClient *bos.Client) *Server {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowWildcard = true
	r.Use(cors.New(config))
	server := &Server{
		s:     r,
		c:     bosClient,
		cache: &RecordCache{},
	}
	r.Use(static.Serve("/", static.LocalFile("frontend/dist", false)))
	r.GET("/api/list", server.getRecordList)
	r.GET("/api/download", server.getDownloadUrl)
	r.GET("/api/status", server.getRestoreStatus)
	r.GET("/api/restore", server.restoreObject)
	r.GET("/api/dsearch", server.search)
	return server
}

func (s *Server) Run() error {
	return s.s.Run("0.0.0.0:8053")
}

type Record struct {
	LiveTitle string    `json:"live_title"`
	StartTime time.Time `json:"start_time"`
	File      []File    `json:"file"`
}

type RecordSlice []*Record

func (a RecordSlice) Len() int           { return len(a) }
func (a RecordSlice) Less(i, j int) bool { return a[i].StartTime.After(a[j].StartTime) }
func (a RecordSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type File struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
	From string `json:"from"`
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func (s *Server) getRecordList(c *gin.Context) {
	user := c.Query("user")
	if user == "" {
		user = "joi"
	}
	s.updateCache()
	c.JSON(200, Response{
		Code: 0,
		Data: s.cache.Records[user],
	})
}

func (s *Server) search(c *gin.Context) {
	room := c.Query("room")
	r, err := strconv.Atoi(room)
	if err != nil {
		c.JSON(200, Response{
			Code: -1,
			Data: err,
		})
		return
	}
	text := c.Query("text")
	c.JSON(200, Response{
		Code: 0,
		Data: db.GetDanmu(r, text),
	})
}

func (s *Server) getDownloadUrl(c *gin.Context) {
	objectName := c.Query("name")
	url := s.c.BasicGeneratePresignedUrl(viper.GetString("bucket"), objectName, 3600) + "&responseContentDisposition=attachment"
	c.JSON(200, Response{
		Code: 0,
		Data: url,
	})
}

func (s *Server) restoreObject(c *gin.Context) {
	objectName := c.Query("name")
	err := s.c.RestoreObject(viper.GetString("bucket"), objectName, 7, api.RESTORE_TIER_STANDARD)
	if err != nil {
		c.JSON(200, Response{
			Code: -1,
			Data: err,
		})
		return
	}
	c.JSON(200, Response{
		Code: 0,
		Data: "success",
	})
}

func (s *Server) getRestoreStatus(c *gin.Context) {
	objectName := c.Query("name")
	meta, err := s.c.GetObjectMeta(viper.GetString("bucket"), objectName)
	if err != nil {
		c.JSON(200, Response{
			Code: -1,
			Data: err,
		})
		return
	}
	if meta.StorageClass != "ARCHIVE" {
		c.JSON(200, Response{
			Code: 0,
			Data: FileStatus{
				File:    objectName,
				Restore: 0,
			},
		})
		return
	}
	if meta.BceRestore == "" {
		c.JSON(200, Response{
			Code: 0,
			Data: FileStatus{
				File:    objectName,
				Restore: 1,
			},
		})
		return
	} else {
		if meta.BceRestore[17] == 't' {
			c.JSON(200, Response{
				Code: 0,
				Data: FileStatus{
					File:    objectName,
					Restore: 2,
				},
			})
			return
		} else {
			c.JSON(200, Response{
				Code: 0,
				Data: FileStatus{
					File:    objectName,
					Restore: 0,
				},
			})
		}
	}
}

// Resore
// 0: no need
// 1: need request
// 2: already going
type FileStatus struct {
	File    string `json:"file"`
	Restore int    `json:"restore"`
}

type RecordCache struct {
	Records    map[string][]*Record
	UpdateTime time.Time
}

func (s *Server) updateCache() {
	if s.cache != nil && (time.Since(s.cache.UpdateTime).Seconds() < 300) && len(s.cache.Records) > 0 {
		return
	}
	s.cache.Records = make(map[string][]*Record)
	s.doupdateCache("joi")
	s.doupdateCache("kiti")
	s.doupdateCache("qilou")
	s.doupdateCache("tocci")
	s.cache.UpdateTime = time.Now()
}
func (s *Server) doupdateCache(user string) {
	log.Info("Start Updating Cache")
	paths := viper.GetStringSlice("paths." + user)
	bucket := viper.GetString("bucket")
	args := new(api.ListObjectsArgs)
	args.Delimiter = "/"
	liveMap := make(map[string]*Record)
	for _, p := range paths {
		args.Prefix = p
		objectResp, err := s.c.ListObjects(bucket, args)
		if err != nil {
			log.Error("Update Cache Failed: ", err)
			return
		}
		log.Infof("Files in %s : %d", p, len(objectResp.Contents))
		for _, entry := range objectResp.Contents {
			if entry.Key[len(entry.Key)-1] == '/' {
				continue
			}
			start, title, err := decodeFileName(entry.Key)
			if err != nil {
				log.Warn(entry.Key)
				continue
			}
			key := fmt.Sprintf("%d-%s", start.YearDay(), title)
			if r, ok := liveMap[key]; ok {
				if r.StartTime.After(start) {
					r.StartTime = start
				}
				r.File = append(r.File, File{
					Name: entry.Key,
					Size: uint64(entry.Size),
					From: p[:2],
				})
			} else {
				liveMap[key] = &Record{
					LiveTitle: title,
					StartTime: start,
				}
				liveMap[key].File = append(liveMap[key].File, File{
					Name: entry.Key,
					Size: uint64(entry.Size),
					From: p[:2],
				})
			}
		}
	}
	log.Infof("Total Live Records: %d", len(liveMap))
	var allRecords []*Record
	for _, v := range liveMap {
		allRecords = append(allRecords, v)
	}
	sort.Sort(RecordSlice(allRecords))
	s.cache.Records[user] = allRecords
}

// [2021-09-07 21-04-02][轴伊Joi_Channel][弹丸论破2].mp4
func decodeFileName(file string) (time.Time, string, error) {
	r1 := regexp.MustCompile(`\[([^\]]+)\]`)
	r2 := regexp.MustCompile(`\](.*)\.(mp4|flv)`)
	infos := r1.FindAllStringSubmatch(file, -1)
	if len(infos) == 0 {
		return time.Now(), "", errors.New("InavlidFileName")
	}
	startTime, err := time.Parse("2006-01-02 15-04-05", infos[0][1])
	if err != nil {
		startTime, err = time.Parse("2006-01-02 15-04", infos[0][1])
		if err != nil {
			return time.Now(), "", err
		}
	}
	var title = ""
	switch len(infos) {
	case 1:
		title = r2.FindAllStringSubmatch(file, -1)[0][1]
	case 2:
		title = r1.FindAllStringSubmatch(file, -1)[1][1]
	case 3:
		title = infos[2][1]
	default:
		return time.Now(), title, errors.New("InavlidFileName")
	}
	return startTime, title, nil
}
