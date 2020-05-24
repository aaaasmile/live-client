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
        { text: 'Name', value: 'Name' },
        { text: 'Date', value: 'DateTimeExp' },
        { text: 'Version', value: 'VersionExp' },
        { text: 'Date Diff', value: 'DateTimeEqual' },
        { text: 'Other Diff', value: 'OtherDiffType' },
        { text: 'Presence', value: 'PresenceType' },
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
    compareDiff(item) {
      console.log('View the item ', item)
      let para = { selected: item.KeyStore, debug: this.debug }
      this.loadingData = true
      API.CompareDiff(this, para)
    },
    getColorPres(prestype, dte) {
      switch (prestype) {
        case "Server only":
          return 'green'
        case "File Source only":
          return 'blue'
        case "Both":
          return 'grey'
      }
    },
    getColorDte(dte) {
      switch (dte) {
        case 'Server is newer':
          return 'purple lighten-2'
        case 'File Source is newer':
          return 'red lighten-3'
      }
      return 'grey'
    },
    getColorOtherDiff(od) {
      switch (od) {
        case 'Diff Modified':
          return 'pink lighten-4'
        case 'Diff Version':
          return 'pink lighten-3'
        case 'Diff Content':
          return 'pink lighten-2'
        case 'Diff Name':
          return 'pink lighten-1'
      }
      return 'grey'
    }
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
      :footer-props="{
      showFirstLastPage: true,
      firstIcon: 'mdi-arrow-collapse-left',
      lastIcon: 'mdi-arrow-collapse-right',
      prevIcon: 'mdi-minus',
      nextIcon: 'mdi-plus'
    }"
    >
      <template v-slot:item.actions="{ item }">
        <v-icon small class="mr-2" @click="compareDiff(item)">mdi-eye</v-icon>
      </template>
      <template v-slot:item.PresenceType="{ item }">
        <v-chip
          :color="getColorPres(item.PresenceType, item.DateTimeEqual)"
          dark
        >{{ item.PresenceType }}</v-chip>
      </template>
      <template v-slot:item.DateTimeEqual="{ item }">
        <v-chip :color="getColorDte(item.DateTimeEqual)" dark>{{ item.DateTimeEqual }}</v-chip>
      </template>
      <template v-slot:item.OtherDiffType="{ item }">
        <v-chip :color="getColorOtherDiff(item.OtherDiffType)" dark>{{ item.OtherDiffType }}</v-chip>
      </template>
    </v-data-table>
  </v-card>
`
}
