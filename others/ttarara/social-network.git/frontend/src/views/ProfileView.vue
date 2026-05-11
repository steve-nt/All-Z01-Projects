<template>
  <div class="page-card">
    <p v-if="profileStore.loading" class="muted">Loading profile...</p>

    <template v-else-if="profileStore.profile && profileStore.profile.limited">
      <ProfileHeader :profile="profileStore.profile" />
      <p class="muted private-notice">
        This profile is private. Follow this user to request full access. You can still
        see any public posts below.
      </p>
      <FollowButton
        v-if="showFollowButton"
        :state="followState"
        :loading="isFollowActionLoading"
        @follow="onRequestFollow"
        @unfollow="onUnfollow"
      />

      <div class="profile-posts-block">
        <h2>Posts</h2>
        <p v-if="profilePostsLoading" class="muted">Loading posts...</p>
        <p v-else-if="profilePostsError" class="error">{{ profilePostsError }}</p>
        <p v-else-if="profilePosts.length === 0" class="muted">
          This user has no posts you can see yet.
        </p>
        <PostList v-else :posts="profilePosts" />
      </div>
    </template>

    <template v-else-if="profileStore.error?.status === 403">
      <h1>Private profile</h1>
      <p class="muted private-notice">
        This profile is private. Follow this user to request full access. You can still
        see any public posts below.
      </p>
      <FollowButton
        v-if="showFollowButton"
        :state="followState"
        :loading="isFollowActionLoading"
        @follow="onRequestFollow"
        @unfollow="onUnfollow"
      />

      <!-- Public or otherwise visible posts for this user -->
      <div class="profile-posts-block">
        <h2>Posts</h2>
        <p v-if="profilePostsLoading" class="muted">Loading posts...</p>
        <p v-else-if="profilePostsError" class="error">{{ profilePostsError }}</p>
        <p v-else-if="profilePosts.length === 0" class="muted">
          This user has no posts you can see yet.
        </p>
        <PostList v-else :posts="profilePosts" />
      </div>
    </template>

    <template v-else-if="profileStore.error">
      <h1>Profile unavailable</h1>
      <p class="error">{{ profileStore.error.message }}</p>
    </template>

    <template v-else-if="profileStore.profile">
      <ProfileHeader :profile="profileStore.profile" />

      <div class="profile-toolbar">
        <template v-if="isOwnProfile">
          <PrivacyToggle
            :is-public="Boolean(profileStore.profile.is_public)"
            :loading="profileStore.updatingPrivacy"
            @change="onTogglePrivacy"
          />
          <button
            type="button"
            class="button secondary settings-trigger"
            @click="settingsModalOpen = true"
          >
            Settings
          </button>
        </template>

        <FollowButton
          v-else-if="showFollowButton"
          :state="followState"
          :loading="isFollowActionLoading"
          @follow="onRequestFollow"
          @unfollow="onUnfollow"
        />
      </div>

      <div class="follow-preview">
        <RouterLink :to="`/profile/${resolvedProfileId}/followers`">
          <strong>{{ followersPrivate ? "Private" : followersCount }}</strong> followers
        </RouterLink>
        <RouterLink :to="`/profile/${resolvedProfileId}/following`">
          <strong>{{ followingPrivate ? "Private" : followingCount }}</strong> following
        </RouterLink>
        <RouterLink v-if="isOwnProfile" to="/follow/requests">
          <strong>{{ pendingRequestsCount }}</strong> follow requests
        </RouterLink>
      </div>

      <!-- Create post (own profile only) -->
      <CreatePostForm
        v-if="isOwnProfile"
        @submitted="onPostSubmitted"
      />

      <!-- Profile posts -->
      <p v-if="profilePostsLoading" class="muted">Loading posts...</p>
      <p v-else-if="profilePostsError" class="error">{{ profilePostsError }}</p>
      <p v-else-if="profilePosts.length === 0" class="muted">No posts yet.</p>
      <PostList v-else :posts="profilePosts" />

      <!-- Settings in modal (opened via Settings button) -->
      <Teleport to="body">
        <div
          v-if="settingsModalOpen"
          class="settings-modal-overlay"
          aria-modal="true"
          role="dialog"
          aria-labelledby="settings-modal-title"
          @click.self="settingsModalOpen = false"
        >
          <div class="settings-modal">
            <div class="settings-modal-header">
              <h2 id="settings-modal-title">Profile settings</h2>
              <button
                type="button"
                class="settings-modal-close"
                aria-label="Close"
                @click="settingsModalOpen = false"
              >
                ×
              </button>
            </div>
            <div class="settings-modal-body">
              <ProfileSettings
                v-if="profileStore.profile"
                :profile="profileStore.profile"
                embedded
              />
            </div>
          </div>
        </div>
      </Teleport>

      <p v-if="localError" class="error">{{ localError }}</p>
      <p v-if="isOwnProfile && postsStore.error" class="error">{{ postsStore.error.message }}</p>
    </template>
  </div>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import { RouterLink, useRoute } from "vue-router";
import CreatePostForm from "../components/CreatePostForm.vue";
import FollowButton from "../components/FollowButton.vue";
import PostList from "../components/PostList.vue";
import PrivacyToggle from "../components/PrivacyToggle.vue";
import ProfileHeader from "../components/ProfileHeader.vue";
import ProfileSettings from "../components/ProfileSettings.vue";
import { socialApi } from "../services/socialApi";
import { useAuthStore } from "../stores/auth";
import { useFollowStore } from "../stores/follow";
import { usePostsStore } from "../stores/posts";
import { useProfileStore } from "../stores/profile";

const route = useRoute();
const authStore = useAuthStore();
const profileStore = useProfileStore();
const followStore = useFollowStore();
const postsStore = usePostsStore();

const localError = ref("");
const profilePosts = ref([]);
const profilePostsLoading = ref(false);
const profilePostsError = ref("");
const followState = ref("follow");
const settingsModalOpen = ref(false);
const followersCount = ref(0);
const followingCount = ref(0);
const pendingRequestsCount = ref(0);
const followersPrivate = ref(false);
const followingPrivate = ref(false);

const resolvedProfileId = computed(() => {
  const idParam = String(route.params.id || "");
  if (idParam === "me") {
    return authStore.userId || 0;
  }
  const parsed = Number(idParam);
  return Number.isInteger(parsed) && parsed > 0 ? parsed : 0;
});

const isOwnProfile = computed(
  () => Boolean(authStore.userId) && authStore.userId === resolvedProfileId.value
);

const showFollowButton = computed(
  () => authStore.loggedIn && !isOwnProfile.value && resolvedProfileId.value > 0
);

const isFollowActionLoading = computed(
  () =>
    followStore.isActionPending(`request:${resolvedProfileId.value}`) ||
    followStore.isActionPending(`unfollow:${resolvedProfileId.value}`)
);

const refreshProfileFollowLists = async () => {
  followersPrivate.value = false;
  followingPrivate.value = false;

  const [followersResult, followingResult] = await Promise.allSettled([
    followStore.fetchFollowers(resolvedProfileId.value),
    followStore.fetchFollowing(resolvedProfileId.value)
  ]);

  if (followersResult.status === "fulfilled") {
    followersCount.value = followersResult.value.length;
  } else if (followersResult.reason?.status === 403) {
    followersPrivate.value = true;
    followersCount.value = 0;
  }

  if (followingResult.status === "fulfilled") {
    followingCount.value = followingResult.value.length;
  } else if (followingResult.reason?.status === 403) {
    followingPrivate.value = true;
    followingCount.value = 0;
  }

  if (isOwnProfile.value && authStore.loggedIn) {
    try {
      const requests = await followStore.fetchPendingRequests();
      pendingRequestsCount.value = requests.length;
    } catch (error) {
      pendingRequestsCount.value = 0;
    }
  } else {
    pendingRequestsCount.value = 0;
  }
};

const resolveFollowState = async () => {
  if (!showFollowButton.value || !authStore.userId) {
    followState.value = "follow";
    return;
  }

  try {
    const data = await socialApi.getFollowing(authStore.userId);
    const myFollowing = data?.following || [];
    const followingMatch = myFollowing.some(
      (item) => Number(item.user_id) === resolvedProfileId.value
    );
    followState.value = followingMatch ? "following" : "follow";
  } catch (error) {
    followState.value = "follow";
  }
};

const loadProfilePage = async () => {
  localError.value = "";
  followState.value = "follow";

  if (!resolvedProfileId.value) {
    localError.value = "Invalid profile id.";
    profileStore.clearProfileState();
    followersCount.value = 0;
    followingCount.value = 0;
    pendingRequestsCount.value = 0;
    followersPrivate.value = false;
    followingPrivate.value = false;
    return;
  }

  try {
    await profileStore.fetchProfile(resolvedProfileId.value);
  } catch (error) {
    if (error?.status !== 403) {
      localError.value = error?.message || "Failed to load profile.";
    }
  }

  await resolveFollowState();

  await refreshProfileFollowLists();

  await loadProfilePosts();
};

const loadProfilePosts = async () => {
  const uid = resolvedProfileId.value;
  if (!uid) {
    profilePosts.value = [];
    return;
  }
  profilePostsLoading.value = true;
  profilePostsError.value = "";
  try {
    const list = await postsStore.fetchPosts({ user_id: uid });
    profilePosts.value = list || [];
    await Promise.all(
      (list || []).map((p) => postsStore.fetchComments(p.post_id))
    );
  } catch (err) {
    profilePosts.value = [];
    profilePostsError.value = err?.message || "Failed to load posts.";
  } finally {
    profilePostsLoading.value = false;
  }
};

const onPostSubmitted = async () => {
  await loadProfilePosts();
};

const onTogglePrivacy = async (isPublic) => {
  localError.value = "";
  try {
    await profileStore.togglePrivacy(isPublic);
  } catch (error) {
    localError.value = error?.message || "Failed to update privacy.";
  }
};

const onRequestFollow = async () => {
  localError.value = "";
  try {
    const response = await followStore.requestFollow(resolvedProfileId.value);
    followState.value = response?.status === "accepted" ? "following" : "requested";
    await resolveFollowState();
    await refreshProfileFollowLists();
  } catch (error) {
    if (error?.status === 400 && error?.data?.error === "duplicate_request") {
      await resolveFollowState();
      if (followState.value !== "following") {
        followState.value = "requested";
      }
      return;
    }
    localError.value = error?.message || "Failed to send follow request.";
  }
};

const onUnfollow = async () => {
  localError.value = "";
  try {
    await followStore.unfollow(resolvedProfileId.value);
    followState.value = "follow";
    await resolveFollowState();
    await refreshProfileFollowLists();
  } catch (error) {
    localError.value = error?.message || "Failed to unfollow.";
  }
};

watch(
  () => route.params.id,
  async () => {
    await loadProfilePage();
  },
  { immediate: true }
);
</script>

<style scoped>
.settings-trigger {
  margin-right: 0;
}

.settings-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.settings-modal {
  background: var(--surface);
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  max-width: 28rem;
  width: 100%;
  max-height: 90vh;
  overflow: auto;
}

.settings-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid var(--border);
}

.settings-modal-header h2 {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
}

.settings-modal-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  line-height: 1;
  color: var(--muted);
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 4px;
}

.settings-modal-close:hover {
  color: var(--surface-text);
  background: var(--border);
}

.settings-modal-body {
  padding: 1.25rem;
}

.settings-modal-body .avatar-upload {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.settings-modal-body .avatar-upload .button {
  flex-shrink: 0;
  cursor: pointer;
}

.private-notice {
  margin-top: 8px;
  margin-bottom: 12px;
}

.profile-posts-block {
  margin-top: 20px;
}

.profile-posts-block h2 {
  font-size: 1.125rem;
  margin: 0 0 12px;
}
</style>
