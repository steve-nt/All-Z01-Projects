import { defineStore } from "pinia";
import { socialApi } from "../services/socialApi";

export const usePostsStore = defineStore("posts", {
  state: () => ({
    posts: [],
    viewerId: null,
    limit: 20,
    offset: 0,
    loading: false,
    creating: false,
    error: null,
    commentsByPostId: {},
    loadingCommentsByPostId: {},
    creatingCommentPostId: null
  }),
  actions: {
    clearError() {
      this.error = null;
    },
    async fetchComments(postId) {
      const id = Number(postId);
      if (!id) return [];
      this.loadingCommentsByPostId = { ...this.loadingCommentsByPostId, [id]: true };
      try {
        const data = await socialApi.getPostComments(id);
        const comments = data?.comments || [];
        this.commentsByPostId = { ...this.commentsByPostId, [id]: comments };
        return comments;
      } catch {
        this.commentsByPostId = { ...this.commentsByPostId, [id]: [] };
        return [];
      } finally {
        this.loadingCommentsByPostId = { ...this.loadingCommentsByPostId, [id]: false };
      }
    },
    async createComment(payload) {
      this.creatingCommentPostId = payload.post_id;
      try {
        await socialApi.createPostComment(payload);
        await this.fetchComments(payload.post_id);
      } catch (error) {
        throw error;
      } finally {
        this.creatingCommentPostId = null;
      }
    },
    async fetchPosts(options = {}) {
      this.loading = true;
      this.clearError();
      try {
        const data = await socialApi.getPosts(options);
        this.posts = data?.posts || [];
        this.viewerId = data?.viewer_id ?? null;
        this.limit = data?.limit ?? this.limit;
        this.offset = data?.offset ?? this.offset;
        return this.posts;
      } catch (error) {
        this.posts = [];
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to load feed"
        };
        throw error;
      } finally {
        this.loading = false;
      }
    },
    async createPost(payload) {
      this.creating = true;
      this.clearError();
      try {
        const result = await socialApi.createPost(payload);
        return result;
      } catch (error) {
        this.error = {
          status: error?.status || 500,
          message: error?.message || "Failed to create post"
        };
        throw error;
      } finally {
        this.creating = false;
      }
    }
  }
});
