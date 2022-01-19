<template>
  <v-app id="inspire">
    <v-navigation-drawer
      app
    >
      <v-sheet
        color="grey lighten-4"
        class="pa-4"
      >
        <v-avatar
          class="mb-4"
          color="orange lighten-3"
          size="64"
        >
          <img src="joi.png" />
        </v-avatar>

        <div>录播管理 @Xinrea</div>
      </v-sheet>

      <v-divider></v-divider>

      <v-list>
        <v-list-item link @click="changeCurrent('joi')">
          <v-list-item-icon>
            <v-icon>mdi-file-video</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>轴伊Joi_Channel</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item link @click="changeCurrent('kiti')">
          <v-list-item-icon>
            <v-icon>mdi-file-video</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>吉吉Kiti</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item link @click="changeCurrent('qilou')">
          <v-list-item-icon>
            <v-icon>mdi-file-video</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>绮楼Qilou</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item link @click="changeCurrent('tocci')">
          <v-list-item-icon>
            <v-icon>mdi-file-video</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>桃星Tocci</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-main>
      <v-container
        class="py-8 px-6"
        fluid
      >
        <file-list :user="current"></file-list>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import FileList from './components/FileList.vue'
  export default {
  components: { FileList },
    data: () => ({
      current: 'joi'
    }),
    methods: {
      changeCurrent(name) {
        this.current = name
      }
    }
  }
</script>