<template>
  <v-flex xs12>
    <v-img
      aspect-ratio="5.75"
      class="teal"
    >
      <v-layout pa-2 column fill-height class="lightbox white--text">
        <v-spacer />
        <v-flex shrink>
          <h2>Project Meow</h2>
          <div>ぞうの卵はおいしいぞう。...</div>
        </v-flex>
      </v-layout>
    </v-img>

    <v-tabs v-model="tab">
      <v-tab ripple>Project</v-tab>
      <v-tab ripple>Comments</v-tab>
      <v-tab ripple>Settings</v-tab>

      <v-tab-item>
        <v-layout row wrap>
          <v-flex xs6>
            <v-container fluid>
              <ogp-card :ogp="{
                image: 'https://avatars2.githubusercontent.com/u/1870091?s=400&amp;v=4',
                title: 'myuon/pyxis',
                description: 'Pyxis. Contribute to myuon/pyxis development by creating an account on GitHub.',
                url: 'https://github.com/myuon/pyxis'
              }"
              />
            </v-container>
          </v-flex>

          <v-flex xs6>
            <v-container fluid>
              <ogp-card :ogp="{
                image: 'https://pbs.twimg.com/media/DtTpA7TU4AAQNx9.jpg:large',
                title: 'みょん on Twitter',
                description: '“様子',
                url: 'https://twitter.com/myuon_myon/status/1068735246793756673'
              }"
              />
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

          <v-layout fluid class="orange lighten-5" style="padding-top: 10px; padding-bottom: 10px;">
            <v-flex shrink style="margin: 10px;">
              <v-avatar color="orange" size="36px">
                <span class="white--text headline">A</span>
              </v-avatar>
            </v-flex>
            <v-flex>
              <div><strong>myuon</strong> - 5分前</div>
              <div>
                へいへいほー<br />
                Yes<br />
                にゃんぽよ
              </div>
            </v-flex>
          </v-layout>

          <v-layout fluid style="padding-top: 10px; padding-bottom: 10px;">
            <v-flex shrink style="margin: 10px;">
              <v-avatar color="orange" size="36px">
                <span class="white--text headline">U</span>
              </v-avatar>
            </v-flex>
            <v-flex>
              <div><strong>anonymous</strong> - 25分前</div>
              <div>
                てすぽよ
              </div>
            </v-flex>
          </v-layout>

          <v-layout fluid style="padding-top: 10px; padding-bottom: 10px;">
            <v-flex shrink style="margin: 10px;">
              <v-avatar color="orange" size="36px">
                <span class="white--text headline">U</span>
              </v-avatar>
            </v-flex>
            <v-flex>
              <div><strong>anonymous</strong> - 25分前</div>
              <div>
                てすぽよ
              </div>
            </v-flex>
          </v-layout>

          <v-layout fluid class="orange lighten-5" style="padding-top: 10px; padding-bottom: 10px;">
            <v-flex shrink style="margin: 10px;">
              <v-avatar color="orange" size="36px">
                <span class="white--text headline">A</span>
              </v-avatar>
            </v-flex>
            <v-flex>
              <div><strong>anonymous</strong> - 45分前</div>
              <div>
                てすぽよ
              </div>
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

export default {
  components: {
    OgpCard,
  },
  data () {
    return {
      tab: 1,
    };
  },
}
</script>
