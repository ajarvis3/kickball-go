import { createRouter, createWebHistory } from "vue-router";

const AtBatPage = () =>
   import("../components/pages/search/AtBatSearchPage.vue");
const GamePage = () => import("../components/pages/search/GameSearchPage.vue");
const PlayerPage = () =>
   import("../components/pages/search/PlayerSearchPage.vue");
const TeamPage = () => import("../components/pages/search/TeamSearchPage.vue");
const LeaguePage = () =>
   import("../components/pages/search/LeagueSearchPage.vue");
const LeagueRulesPage = () =>
   import("../components/pages/search/LeagueRulesSearchPage.vue");
const HelloWorld = () => import("../components/HelloWorld.vue");
const AtBatView = () => import("../components/pages/view/AtBatViewPage.vue");
const GameView = () => import("../components/pages/view/GameViewPage.vue");
const PlayerView = () => import("../components/pages/view/PlayerViewPage.vue");
const TeamView = () => import("../components/pages/view/TeamViewPage.vue");
const LeagueView = () => import("../components/pages/view/LeagueViewPage.vue");
const LeagueRulesView = () =>
   import("../components/pages/view/LeagueRulesViewPage.vue");
const UmpireGameDashboard = () =>
   import("../components/pages/dashboard/UmpireGameDashboard.vue");

const routes = [
   { path: "/", name: "home", component: HelloWorld },
   { path: "/atbats", name: "atbats", component: AtBatPage },
   {
      path: "/atbats/:id",
      name: "atbat-view",
      component: AtBatView,
      props: true,
   },
   { path: "/games", name: "games", component: GamePage },
   { path: "/games/:id", name: "game-view", component: GameView, props: true },
   { path: "/players", name: "players", component: PlayerPage },
   {
      path: "/games/:id/umpire",
      name: "umpire-game",
      component: UmpireGameDashboard,
      props: true,
   },
   { path: "/players/:id",
      name: "player-view",
      component: PlayerView,
      props: true,
   },
   { path: "/teams", name: "teams", component: TeamPage },
   { path: "/teams/:id", name: "team-view", component: TeamView, props: true },
   { path: "/leagues", name: "leagues", component: LeaguePage },
   {
      path: "/leagues/:id",
      name: "league-view",
      component: LeagueView,
      props: true,
   },
   { path: "/leaguerules", name: "leaguerules", component: LeagueRulesPage },
   {
      path: "/leaguerules/:id",
      name: "leaguerules-view",
      component: LeagueRulesView,
      props: true,
   },
];

const router = createRouter({
   history: createWebHistory(),
   routes,
});

export default router;
