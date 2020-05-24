import API from '../apicaller.js'
import TableDiff from '../components/TableDiff.js'

const buildselectedParam = (that) => {
  let arr = []
  that.$store.state.diff.diffSelected.forEach(element => {
    arr.push(element.KeyStore)
  });
  console.log("selection list", arr)
  let para = { selected: arr, debug: that.debug }
  return para
}

export default {
  components: { TableDiff },
  data() {
    return {
      loadingSync: false,
      loadingViewDiff: false,
      loadingExp: false,
      loadingImp: false,
      loadingIgnore: false,
      loadingPage: false,
      samedatefile: false,
      forcesource: false,
      forceserver: false,
      dialogImport: false,
      debug: false,
      selectcount: 500,
      transition: 'scale-transition'
    }
  },
  computed: {
    ...Vuex.mapState({
      Repo: state => {
        return state.gen.prj.repo
      },
      DbLite: state => {
        return state.gen.prj.deblite
      },
    })
  },
  methods: {
    syncRepo() {
      this.loadingSync = true
      let para = {
        forcesource: this.forcesource, 
        forceserver: this.forceserver,
        debug: this.debug
      }
      console.log('Call sync', para)

      API.CallSync(this, para)
    },
    startExplorer(){
      console.log('start explorer TODO...')
    },
    viewDiff() {
      console.log('View diff with beyond compare')
      this.loadingViewDiff = true
      let para = { debug: this.debug }
      API.ViewDiff(this, para)
    },
    importServer() {
      this.dialogImport = false
      let para = buildselectedParam(this)
      this.loadingImp = true
      API.ImportSelectionToServer(this, para)
    },
    exportToFile() {
      let para = buildselectedParam(this)
      this.loadingExp = true
      API.ExportToFile(this, para)
    }
  },
  template: `
  <v-card color="grey lighten-4" flat tile>
    <v-toolbar flat dense>
      <v-toolbar-title class="subheading grey--text">Dashboard</v-toolbar-title>
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="dialogImport = true" :loading="loadingImp" v-on="on">
            <v-icon>mdi-cloud-upload</v-icon>
          </v-btn>
        </template>
        <span>Import into the Server</span>
      </v-tooltip>
      <v-spacer></v-spacer>

      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="viewDiff" :loading="loadingViewDiff" v-on="on">
            <v-icon>mdi-view-comfy</v-icon>
          </v-btn>
        </template>
        <span>View diff</span>
      </v-tooltip>
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="syncRepo" :loading="loadingSync" v-on="on">
            <v-icon>mdi-sync</v-icon>
          </v-btn>
        </template>
        <span>Update view</span>
      </v-tooltip>

      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="exportToFile" :loading="loadingExp" v-on="on">
            <v-icon>mdi-folder</v-icon>
          </v-btn>
        </template>
        <span>Export to local file</span>
      </v-tooltip>
    </v-toolbar>

    <v-row justify="center">
      <v-col class="mb-12" cols="12" md="10">
        <v-skeleton-loader
          :loading="loadingPage"
          :transition="transition"
          height="94"
          type="list-item-two-line"
        >
          <v-card>
            <div class="mx-4">
              <div class="subtitle-2 text--secondary">
                Local Repo: {{Repo}}
                <v-tooltip bottom>
                  <template v-slot:activator="{ on }">
                    <v-btn icon @click="startExplorer" v-on="on">
                      <v-icon>mdi-file</v-icon>
                    </v-btn>
                  </template>
                  <span>View in File Explorer</span>
                </v-tooltip>
              </div>
              <div class="subtitle-2 text--secondary">DB local: {{DbLite}}</div>
            </div>
            <v-expansion-panels :flat="true">
              <v-expansion-panel>
                <v-expansion-panel-header>Options</v-expansion-panel-header>
                <v-expansion-panel-content>
                  <v-row justify="space-around">
                    <v-switch v-model="forcesource" class="ma-2" label="Force Source"></v-switch>
                    <v-switch v-model="forceserver" class="ma-2" label="Force Server Objects"></v-switch>
                    <v-switch v-model="debug" class="ma-2" label="Debug request"></v-switch>
                  </v-row>
                </v-expansion-panel-content>
              </v-expansion-panel>
            </v-expansion-panels>
            <TableDiff></TableDiff>
          </v-card>
        </v-skeleton-loader>
      </v-col>
      <v-row justify="center">
        <v-dialog v-model="dialogImport" persistent max-width="290">
          <v-card>
            <v-card-title class="headline">CAUTION</v-card-title>
            <v-card-text>Do you want to import selected items into the Server?</v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="green darken-1" text @click="importServer">OK</v-btn>
              <v-btn color="green darken-1" text @click="dialogImport = false">Cancel</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-row>
    </v-row>
  </v-card>
`
}