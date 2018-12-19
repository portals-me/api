<template>
  <v-flex xs12>
    <v-img
      aspect-ratio="5.75"
      :class="project.cover && project.cover.color"
    >
      <v-layout pa-2 column fill-height class="lightbox white--text">
        <v-spacer />
        <v-flex shrink>
          <h2>{{ project.title }}</h2>
          <div>{{ project.description }}</div>
        </v-flex>
      </v-layout>
    </v-img>

    <v-tabs v-model="tab">
      <v-tab ripple>Comments</v-tab>
      <v-tab ripple>Articles</v-tab>
      <v-tab ripple>Settings</v-tab>

      <v-tab-item>
        <v-container fluid>
          <v-layout row wrap>
            <v-btn fab small flat>
              <v-icon>add</v-icon>
            </v-btn>
            <v-textarea
              v-model="comment"
              auto-grow
              rows="1"
              single-line
              solo
            />

            <v-btn color="primary" middle @click="submitComment">
              Submit
            </v-btn>
          </v-layout>

          <v-layout fluid :class="comment.owned_by === project.owned_by ? 'orange lighten-5' : ''" style="padding-top: 10px; padding-bottom: 10px;" :key="comment.sort" v-for="comment in project.comments">
            <v-flex shrink style="margin: 10px;">
              <v-avatar color="orange" size="32px">
                <v-img :src="project.members[comment.owned_by].picture" />
              </v-avatar>
            </v-flex>
            <v-flex>
              <div><strong>{{ project.members[comment.owned_by].display_name }}</strong> - {{ new Date(comment.created_at).toISOString() }}</div>
              <div v-html="comment.message.split('\n').join('<br />')"></div>
            </v-flex>
          </v-layout>
        </v-container>
      </v-tab-item>
      <v-tab-item>
        <v-layout row wrap>
          <v-flex xs6 :key="article.id" v-for="article in project.articles">
            <v-container fluid>
              <ogp-card :ogp="article.entity.ogp" />
            </v-container>
          </v-flex>
        </v-layout>
      </v-tab-item>
      <v-tab-item>
      </v-tab-item>
    </v-tabs>
  </v-flex>
</template>

<script>
import OgpCard from '@/components/OgpCard';
import firestore from '@/instance/firestore';
import firebase from 'firebase';
import sdk from '@/app/sdk';

export default {
  components: {
    OgpCard,
  },
  data () {
    return {
      comment: '',
      tab: null,
      project: {
        comments: [],
        articles: [],
      },
      articleIds: [],
    };
  },
  methods: {
    loadArticles () {
      Promise.all(this.articleIds.map(async (articleId, index) => {
        const doc = await firestore.collection('articles').doc(articleId).get();
        this.$set(this.project.articles, index, Object.assign({ id: doc.id }, doc.data()));
      }));
    },
    async loadProject () {
      const projectId = this.$route.params.projectId;
      const project = await sdk.project.get(projectId);
      this.project = Object.assign(project, { articles: [] });
    },
    async submitComment () {
      const projectId = this.$route.params.projectId;
      await sdk.comment.create(projectId, this.comment);

      this.comment = '';
      await this.loadComments();
    },
    async onMount () {
      await this.loadProject();
    },
  },
  mounted: async function () {
    await this.onMount();
  },
}
</script>
