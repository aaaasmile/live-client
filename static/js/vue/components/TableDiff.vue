<template>
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
  </v-card>
</template>