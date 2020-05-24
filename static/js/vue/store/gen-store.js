export default {
  state: {
    errorText: '',
    msgText: '',
    prj: {
      repo: '',
      deblite: '',
    }
  },
  mutations: {
    project(state, prj){
      state.prj.repo = prj.repo
      state.prj.deblite = prj.deblite
    },
    errorText(state, msg) {
      state.errorText = msg
    },
    msgText(state, msg) {
      state.msgText = msg
    },
    clearErrorText(state) {
      if (state.errorText !== '') {
        state.errorText = ''
      }
    },
    clearMsgText(state) {
      if (state.msgText !== '') {
        state.msgText = ''
      }
    }
  }
}