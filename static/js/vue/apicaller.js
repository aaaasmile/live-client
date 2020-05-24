
const handleError = (error, that) => {
  console.error(error);
  that.loadingPage = false
  if (error.bodyText !== '') {
    that.$store.commit('msgText', `${error.statusText}: ${error.bodyText}`)
  } else {
    that.$store.commit('msgText', 'Error: empty response')
  }
}

const handleSetViewDiff = (result, that) => {
  that.loadingPage = false
  console.log('Call terminated ', result.data)
  that.$store.commit('setDiffView', result.data.ResultView)
}

export default {
  CallSync(that, req) {
    that.loadingPage = true
    that.$http.post("CallSync", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
      that.loadingSync = false
      handleSetViewDiff(result, that)
    }, error => {
      that.loadingSync = false
      handleError(error, that)
    });
  },
  ViewDiff(that, req) {
    that.loadingPage = true
    that.$http.post("ViewDiff", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
      that.loadingViewDiff = false
      handleSetViewDiff(result, that)
    }, error => {
      that.loadingViewDiff = false
      handleError(error, that)
    });
  },
  OpenExplorer(that, req) {
    that.$http.post("OpenExplorer", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
      console.log('Open Explorer terminated ', result.data)
      that.$store.commit('msgText', `Status: ${result.data}`)
    }, error => {
      handleError(error, that)
    });
  },
  ExportToFile(that, req) {
    that.loadingPage = true
    that.$http.post("ExportToFile", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
      that.loadingExp = false
      console.log('ExportToFile terminated ', result.data)
      handleSetViewDiff(result, that)
    }, error => {
      that.loadingExp = false
      handleError(error, that)
    });
  },
  NewFile(that, req) {
    that.loadingPage = true
    that.$http.post("NewFile", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
      console.log('NewFile terminated ', result.data)
      handleSetViewDiff(result, that)
    }, error => {
      handleError(error, that)
    });
  },
  ImportSelectionToServer(that, req) {
    that.loadingPage = true
    that.$http.post("ImportSelectionToServer", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
      that.loadingImp = false
      console.log('ImportToNav terminated ', result.data)
      handleSetViewDiff(result, that)
    }, error => {
      that.loadingImp = false
      handleError(error, that)
    });
  },
  CompareDiff(that, req){
    that.$http.post("CompareDiff", JSON.stringify(req), { headers: { "content-type": "application/json" } }).then(result => {
      that.loadingData = false
      console.log('CompareDiff terminated ', result.data)
    }, error => {
      that.loadingData = false
      handleError(error, that)
    });
  }
}