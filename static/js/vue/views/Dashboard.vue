<template>
  <v-card color="grey lighten-4" flat tile>
    <v-toolbar flat dense>
      <v-toolbar-title class="subheading grey--text">Dashboard</v-toolbar-title>
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="dialogImport = true" :loading="loadingImp" v-on="on">
            <v-icon>mdi-cloud-upload</v-icon>
          </v-btn>
        </template>
        <span>Upload</span>
      </v-tooltip>
      <v-spacer></v-spacer>

      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="syncRepo" :loading="loadingSync" v-on="on">
            <v-icon>mdi-sync</v-icon>
          </v-btn>
        </template>
        <span>Update view</span>
      </v-tooltip>
    </v-toolbar>

    <v-row justify="center">
      <v-col class="mb-12" cols="12" md="10">
        <v-card>
          <v-card-title>Live Items</v-card-title>
          <TableDiff></TableDiff>
        </v-card>
      </v-col>
      <v-row justify="center">
        <v-dialog v-model="dialogImport" persistent max-width="290">
          <v-card>
            <v-card-title class="headline">CAUTION</v-card-title>
            <v-card-text>Do you want to import selected items into Live?</v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="green darken-1" text @click="uploadItem">OK</v-btn>
              <v-btn color="green darken-1" text @click="dialogImport = false">Cancel</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-row>
    </v-row>
  </v-card>
</template>