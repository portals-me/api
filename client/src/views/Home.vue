<template>
  <v-flex>
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
              <v-btn flat color="indigo" @click="$router.push(`/collections/${collection.id.split('collection##')[1]}`)">Open</v-btn>
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

<script>
import sdk from '@/app/sdk';

export default {
  data () {
    return {
      dialog: false,
      collections: [],
      form: {
        title: '',
        description: '',
        cover: {
          color: 'teal darken-2',
          sort: 'solid',
        }
      },
    };
  },
  methods: {
    async loadProjects () {
      this.collections = (await sdk.collection.list()).data;
    },
    async createProject () {
      await sdk.collection.create({
        title: this.form.title,
        description: this.form.description,
        cover: this.form.cover,
      });
      this.dialog = false;
      await this.loadProjects();
    },
    async onMount () {
      await this.loadProjects();
    },
  },
  mounted: async function () {
    await this.onMount();
  },
}
</script>
