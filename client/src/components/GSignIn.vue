<template>
  <div id="g-signin2"></div>
</template>

<script>
export default {
  name: 'g-signin',
  methods: {
    onSignIn (user) {
      this.$emit('signin', user);
    },
  },
  mounted () {
    if (typeof gapi !== 'undefined') {
      gapi.signin2.render('g-signin2', {
        onsuccess: this.onSignIn,
      });
    }
  },
}
</script>

