
const handleError = (error, that) => {
  console.error(error);
  if (error.bodyText !== '') {
    that.$store.commit('msgText', `${error.statusText}: ${error.bodyText}`)
  } else {
    that.$store.commit('msgText', 'Error: empty response')
  }
}

export default {
  CallUpload(that, req) {
    that.$http.post("CallUpload", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
    }, error => {
      handleError(error, that)
    });
  },
}