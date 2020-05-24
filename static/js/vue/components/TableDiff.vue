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
</template>