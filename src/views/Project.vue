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
        <v-layout row wrap :key="collection.id" v-for="collection in project.collections">
          <v-flex xs6 :key="item.id" v-for="item in collection.items">
            <v-container fluid>
              <ogp-card :ogp="item.entity.ogp" />
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
      this.project = Object.assign({ id: doc.id }, doc.data());

      const comments = await firestore.collection('projects').doc(projectId).collection('comments').get();
      this.$set(this.project, 'comments', comments.docs.map(doc => {
        return Object.assign({ id: doc.id }, doc.data());
      }));
    }
  },
  mounted: async function () {
    await this.loadProject();
  },
}
</script>
