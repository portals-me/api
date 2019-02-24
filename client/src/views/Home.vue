<template>
  <v-flex>
    <v-subheader>タイムライン</v-subheader>

    <template v-for="(item, index) in timeline">
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
            <a @click="$router.push(`/users/${item.user_name}`)">{{ item.user_display_name }}</a>さんがコレクション<a @click="$router.push(`/collections/${item.item_id.split('collection##')[1]}`)">{{ item.entity.title }}</a>を作りました
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

    <v-subheader>コレクション</v-subheader>

    <v-container grid-list-md>
      <v-layout row wrap>
        <v-flex xs3>
          <v-btn @click="dialog = true" block outline color="indigo" style="margin: 0; height: 100%; min-height: 200px">
            <v-icon left>add</v-icon>
            コレクションを作成
          </v-btn>

          <v-dialog max-width="800" v-model="dialog">
            <v-card>
              <v-card-title class="headline">コレクションを作成</v-card-title>

              <v-card-text>
                <v-form>
                  <v-text-field
                    label="タイトル"
                    v-model="form.title"
                    required
                  />
                  <v-textarea
                    label="説明"
                    v-model="form.description"
                    rows="1"
                    auto-grow
                  />
                </v-form>
              </v-card-text>

              <v-card-actions>
                <v-btn color="success" @click="createProject">
                  Submit
                </v-btn>

                <v-btn flat>
                  Cancel
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-dialog>
        </v-flex>

        <v-flex xs3 v-for="collection in collections" :key="collection.id">
          <v-card>
            <v-img
              aspect-ratio="2.75"
              :class="collection.cover.color"
            >
            </v-img>

            <v-card-title>
              <div>
                <h3 class="headline mb-0">
                  {{ collection.title }}
                </h3>
                <div>{{ collection.description }}</div>
              </div>
            </v-card-title>

            <v-card-actions>
              <v-btn flat color="indigo" @click="$router.push(`/collections/${collection.id}`)">Open</v-btn>
              <v-spacer></v-spacer>
              <v-icon v-if="collection.media && collection.media.includes('document')">edit</v-icon>
              <v-icon v-if="collection.media && collection.media.includes('picture')">brush</v-icon>
              <v-icon v-if="collection.media && collection.media.includes('movie')">movie</v-icon>
            </v-card-actions>
          </v-card>
        </v-flex>
      </v-layout>
    </v-container>
  </v-flex>
</template>

<script lang="ts">
import Vue,{ ComponentOptions } from 'vue';
import { Component } from 'vue-property-decorator';
import VueRouter from 'vue-router';
import Vuex from 'vuex';
import * as types from '@/types';
import sdk from '@/app/sdk';

@Component({
})
export default class Home extends Vue {
  dialog = false;
  form = {
    title: '',
    description: '',
    cover: {
      color: 'teal darken-2',
      sort: 'solid',
    }
  };
  timeline = [];

  public get collections (): Array<types.Collection> {
    return this.$store.state.collections || [];
  }

  async createProject () {
    await sdk.collection.create({
      title: this.form.title,
      description: this.form.description,
      cover: this.form.cover,
    });

    this.dialog = false;
    await this.$store.dispatch('loadCollections', {
      force: true,
    });
  }

  async loadTimeline () {
    const result = await sdk.timeline.get();
    console.log(result);
    this.timeline = result.data;
  }

  async mounted () {
    await Promise.all([
      this.$store.dispatch('loadCollections'),
      this.loadTimeline(),
    ])
  }
}
</script>
