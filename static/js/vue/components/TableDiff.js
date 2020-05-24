// Vuex.mapState in computed (received events)
// Vuex.mapMutations in methods, (emit events)
import API from '../apicaller.js'

export default {
  data() {
    return {
      search: '',
      debug: false,
      loadingData: false,
      headers: [
        { text: 'ID', value: 'ObjectID' },
        { text: 'Name', value: 'Name' },
        { text: 'Date', value: 'Date' },
        { text: 'Version', value: 'Version' },
        { text: 'Actions', value: 'actions', sortable: false },
      ],
    } // end return data()
  },//end data
  computed: {
    diffSelected: {
      get() {
        return (this.$store.state.diff.diffSelected)
      },
      set(newVal) {
        this.$store.commit('setDiffSelected', newVal)
      }
    },
    ...Vuex.mapState({
      diffdata: state => {
        return state.diff.diffdata
      }
    })
  },
  methods: {
    deleteItem(item) {
      console.log('Delete Item TODO... ', item)
      // let para = { selected: item.KeyStore, debug: this.debug }
      // this.loadingData = true
      // API.CompareDiff(this, para)
    },
  },
  template: `
  <v-card>
    <v-card-title>
      <v-text-field v-model="search" append-icon="search" label="Search" single-line hide-details></v-text-field>
    </v-card-title>
    <v-data-table
      v-model="diffSelected"
      :headers="headers"
      :items="diffdata"
      :loading="loadingData"
      item-key="KeyStore"
      show-select
      class="elevation-1"
      :search="search"
    >
      <template v-slot:item.actions="{ item }">
        <v-icon small class="mr-2" @click="deleteItem(item)">mdi-bin</v-icon>
      </template>
    </v-data-table>
  </v-card>`
}
