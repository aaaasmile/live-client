import generic from './gen-store.js'
import diffData from './diffData-store.js'

export default new Vuex.Store({
  modules: {
    gen: generic,
    diff: diffData
  }
})
