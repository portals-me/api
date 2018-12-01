<template>
  <v-flex>
    <v-subheader>Collections</v-subheader>

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

        <v-flex xs3>
          <v-card>
            <v-img
              aspect-ratio="2.75"
              class="teal darken-2"
            >
            </v-img>

            <v-card-title>
              <div>
                <h3 class="headline mb-0">
                  Project Meow
                </h3>
                <div>ぞうの卵はおいしいぞう。ぞうの卵はおいしいぞう。ぞうの卵はおい...</div>
              </div>
            </v-card-title>

            <v-card-actions>
              <v-btn flat color="indigo" @click="$router.push('/project/1')">Open</v-btn>
              <v-spacer></v-spacer>
              <v-icon>edit</v-icon>
              <v-icon>brush</v-icon>
              <v-icon>movie</v-icon>
            </v-card-actions>
          </v-card>
        </v-flex>
      </v-layout>
    </v-container>
  </v-flex>
</template>

<script>
  import HelloWorld from '../components/HelloWorld'

  export default {
    components: {
      HelloWorld
    },
    data () {
      return {
        dialog: false
      };
    }
  }
</script>
