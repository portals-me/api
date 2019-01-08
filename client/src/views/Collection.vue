<template>
  <v-container grid-list-md fluid>
    <v-flex xs12 class="collection-title">
      <h2><v-icon>collections</v-icon> myuon / myuon</h2>
    </v-flex>

    <v-tabs
      class="collection-tab"
    >
      <v-tab>
        <v-icon>layers</v-icon>
        Cards
      </v-tab>

      <v-tab>
        <v-icon>chat</v-icon>
        Channel
      </v-tab>

      <v-tab>
        <v-icon>settings</v-icon>
        Settings
      </v-tab>

      <v-tab-item class="collection-layout">
        <v-layout flex-child wrap>
          <v-flex md3 d-flex>
            <v-sheet
              class="d-flex"
              color="grey lighten-3"
              height="200">
              <sheet-footer>
                Project Meow
              </sheet-footer>
            </v-sheet>
          </v-flex>

          <v-flex md3 d-flex>
            <v-sheet
              class="d-flex"
              color="teal lighten-3"
              height="200">
              <sheet-footer>
                Project Meow
              </sheet-footer>
            </v-sheet>
          </v-flex>

          <v-flex md3 d-flex>
            <v-sheet
              class="d-flex"
              color="red lighten-3"
              height="200">
              <sheet-footer>
                Project Meow
              </sheet-footer>
            </v-sheet>
          </v-flex>

          <v-flex md3 d-flex>
            <v-sheet
              class="d-flex"
              color="indigo lighten-3"
              height="200">
              <sheet-footer>
                Project Meow
              </sheet-footer>
            </v-sheet>
          </v-flex>
        </v-layout>
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
                  <textarea placeholder="myuon/myuonへのメッセージ…" rows="1">
                  </textarea>
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
export default {
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
    }
  }
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

.message textarea {
  border-radius: 3px;
  background-color: rgba(230, 230, 230);
  width: 100%;
  padding: 0.5rem;
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
