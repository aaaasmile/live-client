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
      dialogCreate: false,
      newfilename: '',
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
    newFile(){
      let para = {filename: this.newfilename}
      console.log('New file, param ', para)
      API.NewFile(this, para)
    },
    startExplorer(){
      let para = { repo: this.$store.state.gen.prj.repo }
      console.log('Start Explorer in repo folder', para)
      API.OpenExplorer(this, para)
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
    <v-col class="mb-5" cols="12">
      <v-row justify="center">
        <v-col class="mb-5" cols="12">
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
              <v-card-actions>
                <v-tooltip bottom>
                  <template v-slot:activator="{ on }">
                    <v-btn icon @click="dialogCreate = true" v-on="on">
                      <v-icon>add</v-icon>
                    </v-btn>
                  </template>
                  <span>Create a new File</span>
                </v-tooltip>
              </v-card-actions>
            </v-card>
          </v-skeleton-loader>
        </v-col>
      </v-row>
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
        <v-dialog v-model="dialogCreate" persistent max-width="290">
          <v-card>
            <v-card-title class="headline">New File</v-card-title>
           <v-container>
            <v-col>
              <v-row>
                <v-col cols="10" md="8">
                  <v-text-field v-model="newfilename" label="file name"></v-text-field>
                </v-col>
              </v-row>
            </v-col>
          </v-container>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="green darken-1" text @click="newFile">OK</v-btn>
              <v-btn color="green darken-1" text @click="dialogCreate = false">Cancel</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-row>
    </v-col>
  </v-card>
`
}