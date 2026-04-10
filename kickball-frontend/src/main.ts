import { createApp } from "vue";
import { Notify, Quasar } from "quasar";
import "@quasar/extras/material-icons/material-icons.css";
import "quasar/src/css/index.sass";
import "./style.css";
import App from "./App.vue";
import router from "./router";

createApp(App)
   .use(router)
   .use(Quasar, {
      plugins: { Notify },
      config: {
         brand: {
            primary: "#1976D2",
            secondary: "#26A69A",
            accent: "#9C27B0",
            dark: "#0F1115",
            positive: "#21BA45",
            negative: "#DD3B3B",
            info: "#0EA5E9",
            warning: "#F59E0B",
         },
      },
   })
   .mount("#app");
