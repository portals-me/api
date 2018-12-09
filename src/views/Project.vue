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
      <v-tab ripple>Project</v-tab>
      <v-tab ripple>Comments</v-tab>
      <v-tab ripple>Settings</v-tab>

      <v-tab-item>
        <v-layout row wrap v-if="project.articles">
          <v-flex xs6 :key="article.id" v-for="article in project.articles">
            <v-container fluid>
              <ogp-card :ogp="article.entity.ogp" />
            </v-container>
          </v-flex>
        </v-layout>
      </v-tab-item>
      <v-tab-item>
        <v-container fluid>
          <v-layout row wrap>
            <v-textarea
              auto-grow
              rows="1"
              single-line
              solo
            />

            <v-btn color="primary" middle>
              Submit
            </v-btn>
          </v-layout>

          <v-layout fluid :class="comment.owner.id === project.owner ? 'orange lighten-5' : ''" style="padding-top: 10px; padding-bottom: 10px;" :key="comment.id" v-for="comment in project.comments">
            <v-flex shrink style="margin: 10px;">
              <v-avatar color="orange" size="36px">
                <span class="white--text headline">A</span>
              </v-avatar>
            </v-flex>
            <v-flex>
              <div><strong>{{ comment.owner.user_name }}</strong> - {{ comment.created_at }}</div>
              <div v-html="comment.message.split('\n').join('<br />')"></div>
            </v-flex>
          </v-layout>
        </v-container>
      </v-tab-item>
      <v-tab-item>
      </v-tab-item>
    </v-tabs>
  </v-flex>
</template>

<script>
import OgpCard from '@/components/OgpCard';
import firestore from '@/instance/firestore';

export default {
  components: {
    OgpCard,
  },
  data () {
    return {
      tab: null,
      project: {},
    };
  },
  methods: {
    async loadProject () {
      const projectId = this.$route.params.projectId;
      const doc = await firestore.collection('projects').doc(projectId).get();
      const project = Object.assign({ id: doc.id }, doc.data());

      const articleIds = project.articles;
      const articles = await Promise.all(articleIds.map(async articleId => {
        const doc = await firestore.collection('articles').doc(articleId).get();
        return Object.assign({ id: doc.id }, doc.data());
      }));
      project.articles = articles;

      const comments = await firestore.collection('projects').doc(projectId).collection('comments').get();
      project.comments = comments.docs.map(doc => {
        return Object.assign({ id: doc.id }, doc.data());
      });

      this.project = project;
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
