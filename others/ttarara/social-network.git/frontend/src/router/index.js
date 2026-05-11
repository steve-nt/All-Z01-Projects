import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "../stores/auth";
import { useErrorStore } from "../stores/error";

import HomeView from "../views/HomeView.vue";
import LoginView from "../views/LoginView.vue";
import RegisterView from "../views/RegisterView.vue";
import FeedView from "../views/FeedView.vue";
import ProfileView from "../views/ProfileView.vue";
import GroupsBrowseView from "../views/GroupsBrowseView.vue";
import GroupView from "../views/GroupView.vue";
import GroupInvitationsView from "../views/GroupInvitationsView.vue";
import NotificationsView from "../views/NotificationsView.vue";
import MessagesView from "../views/MessagesView.vue";
import PendingFollowRequestsView from "../views/PendingFollowRequestsView.vue";
import FollowersListView from "../views/FollowersListView.vue";
import FollowingListView from "../views/FollowingListView.vue";
import PeopleView from "../views/PeopleView.vue";

// Part 1: route map (placeholders for future features)
const routes = [
  { path: "/", name: "home", component: HomeView },
  { path: "/login", name: "login", component: LoginView },
  { path: "/register", name: "register", component: RegisterView },
  {
    path: "/feed",
    name: "feed",
    component: FeedView,
    meta: { requiresAuth: true }
  },
  {
    path: "/profile/:id",
    name: "profile",
    component: ProfileView,
    meta: { requiresAuth: true }
  },
  {
    path: "/groups",
    name: "groups",
    component: GroupsBrowseView,
    meta: { requiresAuth: true }
  },
  {
    path: "/groups/:id",
    name: "group",
    component: GroupView,
    meta: { requiresAuth: true }
  },
  {
    path: "/groups/invitations/list",
    name: "group-invitations",
    component: GroupInvitationsView,
    meta: { requiresAuth: true }
  },
  {
    path: "/notifications",
    name: "notifications",
    component: NotificationsView,
    meta: { requiresAuth: true }
  },
  {
    path: "/messages",
    name: "messages",
    component: MessagesView,
    meta: { requiresAuth: true }
  },
  {
    path: "/people",
    name: "people",
    component: PeopleView,
    meta: { requiresAuth: true }
  },
  {
    path: "/follow/requests",
    name: "follow-requests",
    component: PendingFollowRequestsView,
    meta: { requiresAuth: true }
  },
  {
    path: "/profile/:id/followers",
    name: "profile-followers",
    component: FollowersListView,
    meta: { requiresAuth: true }
  },
  {
    path: "/profile/:id/following",
    name: "profile-following",
    component: FollowingListView,
    meta: { requiresAuth: true }
  }
];

// Part 1: SPA router with history mode
const router = createRouter({
  history: createWebHistory(),
  routes
});

// Part 1: auth guards and initial session check
router.beforeEach(async (to) => {
  const authStore = useAuthStore();
  const errorStore = useErrorStore();
  if (errorStore.message) {
    errorStore.clear();
  }
  if (!authStore.sessionChecked && !authStore.isChecking) {
    await authStore.checkSession();
  }

  if (to.meta.requiresAuth && !authStore.loggedIn) {
    return { name: "login", query: { redirect: to.fullPath } };
  }

  if ((to.name === "login" || to.name === "register") && authStore.loggedIn) {
    return { name: "feed" };
  }

  return true;
});

export default router;
