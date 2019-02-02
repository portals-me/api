<template>
  <v-container grid-list-md fluid>
    <v-flex xs12 class="collection-title">
      <v-layout row wrap>
        <v-flex xs5>
          <h2><v-icon>collections</v-icon> {{ collection.owner }} / {{ collection.title }}</h2>
          <pre>{{ collection.description }}</pre>
        </v-flex>
        <v-spacer />
        <v-btn
          dark
          depressed
          outline
          color="indigo"
          @click="createArticleDialog = true"
        >
          <v-icon left>add</v-icon>
          作品を登録
        </v-btn>

        <v-dialog
          v-model="createArticleDialog"
          max-width="760px"
        >
          <v-card>
            <v-card-title
              class="headline grey lighten-2"
              primary-title
            >
              作品を登録
            </v-card-title>

            <v-card-text>
              <v-tabs
                v-model="activeTabInCreateArticleDialog"
                color="indigo lighten-5"
              >
                <v-tab ripple>URL共有</v-tab>
                <v-tab ripple>ファイル</v-tab>

                <v-tab-item>
                  <v-container fluid>
                    <div ref="oEmbedPreview">
                      ここにプレビューが表示されます…
                    </div>
                  </v-container>

                  <v-form
                    v-model="createArticleForm.valid"
                    ref="createArticleForm"
                    lazy-validation
                  >
                    <v-text-field
                      v-model="createArticleForm.url"
                      label="URLを入力"
                      :rules="[v => !!v || '必須項目です']"
                      required
                      @input="previewOEmbed($refs.oEmbedPreview, createArticleForm.url)"
                    />

                    <v-text-field
                      v-model="createArticleForm.title"
                      label="タイトル"
                      :rules="[v => !!v || '必須項目です']"
                      required
                    />

                    <v-textarea
                      v-model="createArticleForm.description"
                      label="説明(任意)"
                      auto-grow
                      rows="1"
                    />

                    <v-checkbox
                      v-model="createArticleForm.checkbox"
                      :rules="[v => !!v || '作品を登録できるのは正当な権利者のみです']"
                      label="私はこの作品の正当な権利者であり、他のいかなる権利の侵害もしていません"
                      required
                    >
                    </v-checkbox>

                    <v-btn
                      :disabled="!createArticleForm.valid"
                      color="indigo"
                      dark
                      @click="submit"
                    >
                      <v-icon left>send</v-icon>
                      送信
                    </v-btn>
                  </v-form>
                </v-tab-item>

                <v-tab-item>
                  <v-form
                    v-model="createArticleForm.valid"
                    ref="createArticleForm"
                    lazy-validation
                  >
                    <input
                      type="file"
                      multiple
                      accept="image/*"
                      @change="selectFiles"
                    />
                    <div class="preview" v-if="imageData">
                      <img
                        v-for="(file, index) in imageData"
                        :key="index"
                        :src="file.src"
                      >
                    </div>

                    <v-text-field
                      v-model="createArticleForm.title"
                      label="タイトル"
                      :rules="[v => !!v || '必須項目です']"
                      required
                    />

                    <v-textarea
                      v-model="createArticleForm.description"
                      label="説明(任意)"
                      auto-grow
                      rows="1"
                    />

                    <v-checkbox
                      v-model="createArticleForm.checkbox"
                      :rules="[v => !!v || '作品を登録できるのは正当な権利者のみです']"
                      label="私はこの作品の正当な権利者であり、他のいかなる権利の侵害もしていません"
                      required
                    >
                    </v-checkbox>

                    <v-btn
                      :disabled="!createArticleForm.valid"
                      color="indigo"
                      dark
                      @click="submit"
                    >
                      <v-icon left>send</v-icon>
                      送信
                    </v-btn>
                  </v-form>
                </v-tab-item>
              </v-tabs>

            </v-card-text>
          </v-card>
        </v-dialog>
      </v-layout>
    </v-flex>

    <v-tabs
      class="collection-tab"
    >
      <v-tab>
        <v-icon>layers</v-icon>
        作品
      </v-tab>

      <v-tab>
        <v-icon>chat</v-icon>
        チャンネル
      </v-tab>

      <v-tab>
        <v-icon>settings</v-icon>
        設定
      </v-tab>

      <v-tab-item class="collection-layout">
        <v-layout flex-child wrap>
          <v-hover :key="'v-hover-' + index" v-for="(article, index) in articles">
            <v-flex
              md3
              d-flex
              slot-scope="{ hover }"
              @click="clickArticleCard(index)"
            >
              <v-card
                class="mx-auto"
                :class="`elevation-${hover ? 6 : 2}`"
              >
                <v-img
                  :aspect-ratio="16/9"
                  :src="article.entity.type == 'image' ? article.entity.url : 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVQYV2N4+PDhfwAI7QOjRSIQaQAAAABJRU5ErkJggg=='"
                />
                <v-card-title>
                  {{ article.title }}
                </v-card-title>
              </v-card>
            </v-flex>
          </v-hover>
        </v-layout>

        <v-dialog
          v-model="articleDialog"
          max-width="600"
        >
          <v-card>
            <v-card-title class="headline" v-if="activeArticle.id">{{ activeArticle.title }}</v-card-title>

            <v-card-text>
              <pre v-if="activeArticle.id">{{ activeArticle.description }}</pre>

              <p v-if="activeArticle.id && activeArticle.entity.type == 'share'"><a :href="activeArticle.entity.url">{{ activeArticle.entity.url }}</a></p>

              <v-img
                :src="activeArticle.entity.url"
                v-if="activeArticle.id && activeArticle.entity.type == 'image'"
              />
              <div ref="articleDialog"></div>
            </v-card-text>
          </v-card>
        </v-dialog>
      </v-tab-item>

      <v-tab-item>
        <v-layout>
          <v-flex xs12 class="message">
            <v-avatar>
              <v-img
                src="https://lh6.googleusercontent.com/-HrqEjsNu_No/AAAAAAAAAAI/AAAAAAAAAMI/Rg4RwE9Y7So/s96-c/photo.jpg"
              />
            </v-avatar>

            <div class="content">
              <div class="header">
                <strong>myuon</strong>
              </div>

              <div class="input-area">
                <div>
                  <autogrow-textarea placeholder="myuon/myuonへのメッセージ…" />
                </div>

                <v-btn depressed color="primary">
                  送信
                </v-btn>
              </div>
            </div>
          </v-flex>
        </v-layout>

        <v-divider />

        <v-layout>
          <v-flex xs12 class="message">
            <v-avatar>
              <v-img
                src="https://lh6.googleusercontent.com/-HrqEjsNu_No/AAAAAAAAAAI/AAAAAAAAAMI/Rg4RwE9Y7So/s96-c/photo.jpg"
              />
            </v-avatar>

            <div class="content">
              <div class="header">
                <strong>myuon</strong>

                10分前
              </div>

              <div class="content-content">
                This is 最高にちょうどいい本文。
              </div>
            </div>
          </v-flex>
        </v-layout>

        <v-layout>
          <v-flex xs12 class="message">
            <v-avatar>
              <v-img
                src="https://lh6.googleusercontent.com/-HrqEjsNu_No/AAAAAAAAAAI/AAAAAAAAAMI/Rg4RwE9Y7So/s96-c/photo.jpg"
              />
            </v-avatar>

            <div class="content">
              <div class="header">
                <strong>myuon</strong>

                10分前
              </div>

              <div class="content-content">
                This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。This is 最高にちょうどいい本文。
              </div>
            </div>
          </v-flex>
        </v-layout>

        <v-layout>
          <v-flex xs12 class="message">
            <v-avatar>
              <v-img
                src="https://lh6.googleusercontent.com/-HrqEjsNu_No/AAAAAAAAAAI/AAAAAAAAAMI/Rg4RwE9Y7So/s96-c/photo.jpg"
              />
            </v-avatar>

            <div class="content">
              <div class="header">
                <strong>myuon</strong>

                10分前
              </div>

              <div class="content-content">
                L <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                o <br />
                g cat <br />
              </div>
            </div>
          </v-flex>
        </v-layout>
      </v-tab-item>

      <v-tab-item>
        <v-container>
          <v-btn
            color="error"
            @click="deleteCollection"
          >
            プロジェクトを削除
            <v-icon right dark>delete</v-icon>
          </v-btn>
        </v-container>
      </v-tab-item>
    </v-tabs>
  </v-container>
</template>

<script>
import AutogrowTextarea from '@/components/AutogrowTextarea';
import fetchJsonp from 'fetch-jsonp';
import sdk from '@/app/sdk';
import axios from 'axios';

export default {
  data () {
    return {
      activeTabInCreateArticleDialog: null,
      articleDialog: false,
      createArticleDialog: false,
      createArticleDialogTab: null,
      createArticleForm: {
        valid: true,
        url: '',
        checkbox: false,
        title: '',
        description: '',
      },
      collection: {},
      articles: [],
      activeArticle: {},
      imageData: [],
    };
  },
  components: {
    SheetFooter: {
      functional: true,

      render (h, { children }) {
        return h('v-sheet', {
          staticClass: 'mt-auto align-center justify-center d-flex',
          props: {
            color: 'rgba(0, 0, 0, .36)',
            dark: true,
            height: 50
          }
        }, children)
      }
    },
    AutogrowTextarea,
  },
  methods: {
    selectFiles (input) {
      Array.from(input.target.files).forEach((file, index) => {
        const url = URL.createObjectURL(file);
        this.imageData.push({ src: url, file: file });
      });
    },
    clickArticleCard (index) {
      this.articleDialog = true;
      this.activeArticle = this.articles[index];

      if (this.activeArticle.entity.type === 'share') {
        this.previewOEmbed(this.$refs.articleDialog, this.activeArticle.entity.url);
      } else {
        this.$refs.articleDialog.innerHTML = ``;
      }
    },
    async previewOEmbed (elem, url_raw) {
      const getProvider = (url) => {
        if (/https:\/\/twitter\.com\/.*\/status\/.*/.test(url)) {
          return `https://publish.twitter.com/oembed?format=json&url=${encodeURIComponent(url)}`
        }
      };

      const url = getProvider(url_raw);

      if (!url) {
        elem.innerText = 'ここにプレビューが表示されます…';
        return;
      };

      const response = await fetchJsonp(url);
      const card_json = await response.json();

      const replaceHTML = (element, html) => {
        element.innerHTML = html;
        element.querySelectorAll('script').forEach(scriptElement => {
          const se = document.createElement('script');
          se.src = scriptElement.src;
          scriptElement.replaceWith(se);
        });
      }

      replaceHTML(elem, card_json.html);
    },
    async loadArticles () {
      const collectionId = this.$route.params.collectionId;
      const result = (await sdk.article.list(collectionId)).data;
      this.articles = result;
    },
    async submit () {
      if (this.$refs.createArticleForm.validate()) {
        await this.postArticle();
      }
    },
    async postArticle () {
      const collectionId = this.$route.params.collectionId;

      if (this.imageData.length != 0) {
        const presignedURL = (await sdk.article.generate_presigned_url(collectionId, this.imageData[0].file.name)).data;
        await axios.put(presignedURL, this.imageData[0].file, {
          headers: { 'Content-Type': this.imageData[0].file.type },
        });

        const user = JSON.parse(localStorage.getItem('user'));
        await sdk.article.create(collectionId, {
          title: this.createArticleForm.title,
          description: this.createArticleForm.description,
          entity: {
            type: "image",
            format: "png",
            url: `https://s3-ap-northeast-1.amazonaws.com/portals-me-storage-users/${encodeURIComponent(user.id)}/${collectionId}/${this.imageData[0].file.name}`,
          }
        });
      } else {
        await sdk.article.create(collectionId, {
          title: this.createArticleForm.title,
          description: this.createArticleForm.description,
          entity: {
            type: "share",
            format: "oembed",
            url: this.createArticleForm.url,
          }
        });
      }

      this.createArticleDialog = false;
      await this.loadArticles();
    },
    async loadCollection () {
      const collectionId = this.$route.params.collectionId;
      const collection = (await sdk.collection.get(collectionId)).data;
      this.collection = Object.assign(collection);
    },
    async deleteCollection () {
      const collectionId = this.$route.params.collectionId;
      await sdk.collection.delete(collectionId);
      this.$router.push('/dashboard');
    },
  },
  async mounted () {
    await Promise.all([
      this.loadCollection(),
      this.loadArticles(),
    ]);
  },
}
</script>

<style scoped>
.collection-title {
  margin-bottom: 1em;
}

.collection-tab .v-icon {
  margin-right: 0.2em;
}

.collection-layout .card {
  margin-top: 1rem;
}

.message {
  display: flex;
}

.message .v-avatar {
  display: block;
  flex: 0 0 auto;
  height: auto;
  margin-right: 20px;
}

.message .content {
  display: block;
  flex: 1 1 auto;
  margin-top: 1em;
}

.message .input-area {
  display: flex;
}

.message .input-area > div {
  display: block;
  flex: 1 1 auto;
}

.message .input-area > .v-btn {
  margin-top: 0;
  margin-bottom: 0;
}

.message .content .header {
  margin-bottom: 0.3em;
}
</style>
