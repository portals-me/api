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
                New Project New Project New Project New Project New Project
              </v-card-text>

              <v-card-actions>
                <v-btn color="success">
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
              <v-icon v-if="project.media.includes('document')">edit</v-icon>
              <v-icon v-if="project.media.includes('picture')">brush</v-icon>
              <v-icon v-if="project.media.includes('movie')">movie</v-icon>
            </v-card-actions>
          </v-card>
        </v-flex>
      </v-layout>
    </v-container>
  </v-flex>
</template>

<script>
import sdk from '@/sdk';

export default {
  data () {
    return {
      dialog: false,
      projects: [],
    };
  },
  methods: {
    async loadProjects () {
      this.projects = await sdk(this.$config.API, this.$config.axios).project.list();
    },
  },
  mounted: async function () {
    await this.loadProjects();
  },
}
</script>
