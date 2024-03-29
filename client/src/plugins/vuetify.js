/**
 * plugins/vuetify.js
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import "@mdi/font/css/materialdesignicons.css";
import "vuetify/styles";

import { createVuetify } from "vuetify";
import { VDataTable } from "vuetify/components/VDataTable";

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    themes: {
      light: {
        colors: {
          primary: "#217dbb",
          secondary: "#5CBBF6",
          background: "#D8F2FA",
          accent: "#FFFFFF",
        },
      },
      dark: {
        colors: {
          primary: "#217dbb",
          secondary: "#5CBBF6",
          background: "#333",
        },
      },
    },
  },
  components: {
    VDataTable,
  },
  defaults: {
    VDataTable: {
      fixedHeader: true,
      noDataText: "Results not found",
    },
  },
});
