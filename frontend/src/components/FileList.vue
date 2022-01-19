<template>
  <v-row>
    <v-col
      cols="6"
    >
      <v-card>
        <v-card-title>{{getFullName(this.user)}} 录播列表</v-card-title>
        <v-list two-line>
            <v-list-item-group>
              <v-list-item
                v-for="item in files"
                :key="item.live_title+item.start_time"
                @click="showDetail(item)"
              >
                <v-list-item-content>
                  <v-list-item-title>{{ item.live_title }}</v-list-item-title>
                  <v-list-item-subtitle>
                    {{ item.start_time.replace('T',' ').replace('Z','') }}
                  </v-list-item-subtitle>
                </v-list-item-content>
              </v-list-item>
            </v-list-item-group>
        </v-list>
      </v-card>
    </v-col>
    <v-col
      cols="6"
      >
      <v-card style="margin-right: 20px">
        <v-card-title>弹幕搜索</v-card-title>
        <v-card-text>
          <v-text-field v-model="searchText"/>
          <v-btn @click="search">Search</v-btn>
            <v-data-table
            :headers="table.headers"
            :items="table.danmus"
            :items-per-page="10"
          ></v-data-table>
        </v-card-text>
      </v-card>
    </v-col>
    <v-dialog
      v-model="dialog"
    >
      <v-card>
        <v-card-title>
          {{ showedItem.live_title }}
        </v-card-title>
        <v-card-text>
          <v-list>
            <v-list-item-group>
              <v-list-item
                v-for="f in showedItem.file"
                :key="f.name"
              >
                <v-chip 
                  class="mr-3"
                  :color="getColorFrom(f.from)"
                  style="color: white"
                  small
                >
                  {{ f.from }}
                </v-chip>
                <v-chip 
                  class="mr-3"
                  color="grey"
                  style="color: white"
                  small
                >
                  {{ formatSize(f.size) }}
                </v-chip>
                {{ getFileName(f.name) }}
                <v-spacer />
                <v-list-item-action>
                  <v-btn plain
                    v-if="f.valid == 0"
                    @click="download(f.name)"
                  >
                    <v-icon>
                      mdi-download
                    </v-icon>
                    下载
                  </v-btn>
                  <v-btn plain
                    v-if="f.valid == 1"
                    @click="restore(f.name)"
                  >
                    <v-icon>
                      mdi-reload
                    </v-icon>
                    取回
                  </v-btn>
                  <div v-if="f.valid == 2" style="color: green">
                    正在取回
                  </div>
                </v-list-item-action>
              </v-list-item>
            </v-list-item-group>
          </v-list>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<script>
import axios from 'axios'
import '../plugins/format'
const endpoint = "https://record.joi-club.cn"
// const endpoint = "http://localhost:8053"
  export default {
    name: 'FileList',
    props: ['user'],
    data() {
      return {
        files: [],
        dialog: false,
        showedItem: {
          live_title: ""
        },
        searchText: "",
        table: {
          headers: [
            { text: '时间', value: 'create_time' },
            { text: '用户', value: 'user_name' },
            { text: '内容', value: 'content' },
          ],
          danmus: []
        }
      }
    },
    watch: {
      user: function(val) {
        console.log(val)
        axios.get(endpoint+"/api/list?user="+this.user).then(resp=>{
          let data = resp.data
          if (data.code === 0) {
            this.files = data.data
          }
        })
      }
    },
    mounted() {
      axios.get(endpoint+"/api/list?user="+this.user).then(resp=>{
        let data = resp.data
        if (data.code === 0) {
          this.files = data.data
        }
      })
    },
    methods: {
      search() {
        if (this.searchText === "") {
          return
        }
        console.log(this.searchText)
        axios.get(endpoint+"/api/dsearch?room="+this.getRoom(this.user)+"&text="+this.searchText).then(resp=>{
          let data = resp.data
          if (data.code === 0) {
            if (data.data === null) {
              this.table.danmus = []
            } else {
              data.data.forEach(item=>{
                let t = new Date(item.create_time)
                item.create_time = t.format('yyyy-mm-dd HH:MM:ss')
              })
              this.table.danmus = data.data
            }
          }
        })
      },
      getFileName(name) {
        return name.split('/', 3)[2]
      },
      getFullName(name) {
        switch (name) {
          case 'joi':
            return '轴伊Joi_Channel'
          case 'kiti':
            return '吉吉Kiti'
          case 'qilou':
            return '绮楼Qilou'
          case 'tocci':
            return '桃星Tocci'
          default:
            return ''
        }
      },
      getRoom(name) {
        switch (name) {
          case 'joi':
            return 21484828
          case 'kiti':
            return 23017343
          case 'qilou':
            return 23017346
          case 'tocci':
            return 23017349
          default:
            return 21484828
        }
      },
      showDetail(item) {
        this.showedItem = item
        this.updateShowedItem()
        this.dialog = true
      },
      updateShowedItem() {
        this.showedItem.file.forEach((f,i,a) => {
          axios.get(endpoint+"/api/status?name="+f.name).then(resp=>{
            let res = resp.data
            if (res.code === 0) {
              a[i].valid = res.data.restore
            } else {
              a[i].valid = 0
            }
            this.showedItem.file = [...this.showedItem.file]
          })
        });
      },
      getColorFrom(from) {
        switch (from) {
          case "S1":
            return "orange"
          case "S2":
            return "blue"
          case "S3":
            return "green"
          default:
            return "grey"
        }
      },
      download(name) {
        axios.get(endpoint+"/api/download?name="+name).then(resp=>{
          let data = resp.data
          if (data.code === 0) {
            window.open(data.data)
          } else {
            console.error(data)
          }
        })
      },
      restore(name) {
        axios.get(endpoint+"/api/restore?name="+name).then(resp=>{
          let data = resp.data
          if (data.code === 0) {
            this.updateShowedItem()
          } else {
            console.error(data)
          }
        })
      },
      formatSize(size) {
        let suffix = ' Bytes'
        let value = size
        if (value > 1024) {
          value = value/1024
          suffix = ' KiB'
        } 
        if (value > 1024) {
          value = value/1024
          suffix = ' MiB'
        } 
        if (value > 1024) {
          value = value/1024
          suffix = ' GiB'
        } 
        return value.toFixed(2)+suffix
      }
    }
  }
</script>
