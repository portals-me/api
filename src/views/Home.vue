<template>
  <v-flex>
    <v-subheader>Projects</v-subheader>

    <v-container grid-list-md>
      <v-layout row wrap>
        <v-flex xs3>
          <v-btn @click="dialog = true" block outline color="indigo" style="margin: 0; height: 100%;">
            <v-icon left>add</v-icon>
            New Project
          </v-btn>

          <v-dialog max-width="800" v-model="dialog">
            <v-card>
              <v-card-title class="headline">Create A New Project</v-card-title>

              <v-card-text>
                <v-form>
                  <v-text-field
                    label="Title"
                    v-model="form.title"
                    required
                  />
                  <v-textarea
                    label="Project Description"
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

        <v-flex xs3 v-for="project in projects" :key="project.id">
          <v-card>
            <v-img
              aspect-ratio="2.75"
              :class="project.cover.color"
            >
            </v-img>

            <v-card-title>
              <div>
                <h3 class="headline mb-0">
                  {{ project.title }}
                </h3>
                <div>{{ project.description }}</div>
              </div>
            </v-card-title>

            <v-card-actions>
              <v-btn flat color="indigo" @click="$router.push(`/projects/${project.id}`)">Open</v-btn>
              <v-spacer></v-spacer>
              <v-icon v-if="project.media && project.media.includes('document')">edit</v-icon>
              <v-icon v-if="project.media && project.media.includes('picture')">brush</v-icon>
              <v-icon v-if="project.media && project.media.includes('movie')">movie</v-icon>
            </v-card-actions>
          </v-card>
        </v-flex>
      </v-layout>
    </v-container>
  </v-flex>
</template>

<script>
import firestore from '@/instance/firestore';
import firebase from 'firebase';

export default {
  data () {
    return {
      dialog: false,
      projects: [],
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
      const projects = await firestore
        .collection('projects')
        .where('owner', '==', this.$store.state.user.uid)
        .orderBy('created_at', 'desc')
        .get();
      this.projects = projects.docs.map(doc => {
        return Object.assign({ id: doc.id }, doc.data());
      });
    },
    async createProject () {
      const project = {
        title: this.form.title,
        owner: this.$store.state.user.uid,
        cover: this.form.cover,
        created_at: firebase.firestore.FieldValue.serverTimestamp(),
      };
      await firestore
        .collection('projects')
        .add(project);
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
