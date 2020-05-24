export default {
  state: {
    diffSelected: [],
    diffdata: [
    ],
  },
  mutations: {
    setDiffSelected(state, selected) {
      state.diffSelected = selected
    },
    setDiffView(state, dataView) {
      state.diffSelected = []
      state.diffdata = dataView
    },
    selectAll(state,count) {
      state.diffSelected = state.diffdata.slice(0,count)
    }
  }
}