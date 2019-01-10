<template>
  <textarea :placeholder="$props.placeholder" rows="1" @keydown="autosize" ref="textarea"></textarea>
</template>

<script>
export default {
  props: [ 'placeholder' ],
  methods: {
    autosize (event) {
      const elem = this.$refs.textarea;
      setTimeout(() => {
        elem.style.height = 0;
        elem.style.height = `${Math.max(14, elem.scrollHeight)}px`;
      }, 0);
    }
  },
}
</script>

<style scoped>
textarea {
  overflow: hidden;
  padding: 0.5rem;
  width: 100%;
  font-size: 14px;
  display: block;
  border-radius: 3px;
  background-color: rgba(230, 230, 230);
}
</style>
