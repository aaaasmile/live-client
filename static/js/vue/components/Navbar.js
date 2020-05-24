import routes from '../routes.js'
import Toast from './toast.js'

export default {
  components: {Toast},
  data() {
    return {
      drawer: false,
      links: routes,
      AppTitle: "Client"
    }
  },
  template: `
  <nav>
    <v-app-bar dense flat>
      <v-btn text color="grey">
        <v-icon>mdi-menu</v-icon>
      </v-btn>
      <v-toolbar-title class="text-uppercase grey--text">
        <span class="font-weight-light">Live</span>
        <span>{{AppTitle}}</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>
    </v-app-bar>
    <Toast></Toast>
  </nav>
`
}