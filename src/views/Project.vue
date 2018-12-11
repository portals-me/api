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

          <v-layout fluid :class="comment.owner.id === project.owner ? 'orange lighten-5' : ''" style="padding-top: 10px; padding-bottom: 10px;" :key="comment.id" v-for="comment in project.comments">
            <v-flex shrink style="margin: 10px;">
              <v-avatar color="orange" size="36px">
                <span class="white--text headline">A</span>
              </v-avatar>
            </v-flex>
            <v-flex>
              <div><strong>{{ comment.owner.display_name }}</strong> - {{ comment.created_at.toDate().toISOString() }}</div>
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
    async loadComments () {
      const projectId = this.$route.params.projectId;
      const comments = await firestore.collection('projects').doc(projectId).collection('comments').orderBy('created_at', 'desc').limit(30).get();
      await Promise.all(comments.docs.map(async (doc, index) => {
        const comment = doc.data();
        const user = await firestore.collection('users').doc(comment.owner).get();

        this.$set(this.project.comments, index, Object.assign(comment, { owner: Object.assign(user.data(), { id: comment.owner })}));
      }));
    },
    loadArticles () {
      Promise.all(this.articleIds.map(async (articleId, index) => {
        const doc = await firestore.collection('articles').doc(articleId).get();
        this.$set(this.project.articles, index, Object.assign({ id: doc.id }, doc.data()));
      }));
    },
    async loadProject () {
      const projectId = this.$route.params.projectId;
      const doc = await firestore.collection('projects').doc(projectId).get();
      this.articleIds = doc.data().articles;
      this.project = Object.assign(doc.data(), { id: doc.id, comments: [], articles: [] });

      await Promise.all([
        this.loadArticles(),
        this.loadComments(),
      ]);
    },
    async submitComment () {
      const projectId = this.$route.params.projectId;
      await firestore.collection('projects').doc(projectId).collection('comments').add({
        owner: this.$store.state.user.uid,
        message: this.comment,
        created_at: firebase.firestore.FieldValue.serverTimestamp(),
      });

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
