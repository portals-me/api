<template>
  <v-container>
    <v-layout justify-center>
      <v-flex xs8>
        <v-avatar
          size="120"
          color="grey lighten-4"
        >
          <img :src="user.picture" alt="avatar">
        </v-avatar>

        <h2>{{ user.display_name }}</h2>
        <p>@{{ user.name }}</p>

        <v-dialog
          v-model="editDialog"
          width="500"
          v-if="user.id == me.id"
        >
          <v-btn
            outline
            color="indigo"
            style="margin-left: 0;"
            slot="activator"
          >
            プロフィールを編集
          </v-btn>

          <v-card>
            <v-card-title>
              <edit-user-basic-profile
                :formData="{ name: user.name, display_name: user.display_name, picture: user.picture }"
                @submit="updateUserProfile"
              />
            </v-card-title>
          </v-card>
        </v-dialog>

        <v-btn
          outline
          color="indigo"
          style="margin-left: 0;"
          @click="follow"
          v-if="user.id != me.id"
        >
          このユーザーをフォロー
        </v-btn>

        <v-flex xs12>
          <v-list>
            ユーザーアクティビティ

            <template v-for="(item, index) in feed">
              <v-divider
                :key="'div-before-' + index"
                v-if="index == 0"
              />

              <v-list-tile
                :key="index"
                two-line
              >
                <v-list-tile-content>
                  <v-list-tile-title v-if="item.event_name == 'INSERT_COLLECTION'">
                    コレクション<a @click="$router.push(`/collections/${item.item_id.split('collection##')[1]}`)">{{ item.entity.title }}</a>を作りました
                  </v-list-tile-title>
                  <v-list-tile-title v-if="item.event_name == 'INSERT_ARTICLE'">
                    作品<a @click="$router.push(`/collections/${item.item_id.split('/')[0].split('collection##')[1]}`)">{{ item.entity.title }}</a>を投稿しました
                  </v-list-tile-title>
                  <v-list-tile-sub-title v-if="item.entity.description">{{ item.entity.description }}</v-list-tile-sub-title>
                </v-list-tile-content>

                <v-list-tile-action>
                  <v-list-tile-action-text>{{ new Date(item.timestamp * 1000).toLocaleString() }}</v-list-tile-action-text>
                </v-list-tile-action>
              </v-list-tile>

              <v-divider
                :key="'div-' + index"
              />
            </template>
          </v-list>
        </v-flex>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script lang="ts">
import Vue,{ ComponentOptions } from 'vue';
import { Component } from 'vue-property-decorator';
import VueRouter from 'vue-router';
import sdk from '@/app/sdk';
import EditUserBasicProfile from '@/components/EditUserBasicProfile.vue';

@Component({
  components: {
    EditUserBasicProfile,
  }
})
export default class User extends Vue {
  user: any = null;
  feed: Array<any> = [];
  me: object = {};
  editDialog = false;

  async follow () {
    const userName = this.$route.params.userId;
    await sdk.user.follow(userName);
  }

  async updateUserProfile (form: any) {
    console.log(form);
  }

  async mounted () {
    await Promise.all([
      (async () => {
        const userName = this.$route.params.userId;
        const result = await sdk.user.get(userName);
        this.user = result.data;
      })(),
      (async () => {
        const userName = this.$route.params.userId;
        const result = await sdk.user.feed.list(userName);
        this.feed = result.data;
      })(),
    ]);

    try {
      this.me = JSON.parse(localStorage.getItem('user') as string);
    } catch (e) {
    }
  }
}
</script>
