// ============================================================================
// SECTION: FILE ROLE & INVARIANTS
// ============================================================================
// Composition root: wires global state, lazy-loads feature modules, initializes routing/auth/WebSocket, and exposes globals used by inline handlers.
// Must NOT own feature logic (delegated to modules) or change DOM structure/IDs.
// Invariants: preserve init order (auth → WebSocket → router → event listeners → status box drag), keep AppState shape stable, and maintain global function/ID contracts.

// ============================================================================
// SECTION: GLOBAL CONSTANTS & SHARED STATE
// ============================================================================
// Global state
const AppState = {
    currentUser: null,
    currentPage: 'home',
    ws: null,
    currentChatUser: null,
    messagesOffset: 0,
    isLoadingMessages: false,
    hasMoreMessages: true,
    unreadMessages: {}, // Track unread message counts per user ID
    unreadNotifications: 0, // Track unread notification count
    notifications: [], // Store notifications
    typingTimeouts: {}, // Store timeouts for typing indicators
    typingDebounceTime: 3000, // 3 seconds delay before hiding
    chatScrollHandler: null, // Store scroll handler for cleanup
    updateLoadMoreButton: null, // Function to update load more button visibility
};

// ============================================================================
// SECTION: APPLICATION STARTUP (DOM READY & VISIBILITY)
// ============================================================================
// Initialize app
document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM loaded, initializing app...');

    // DOM contract: users status box must exist; hidden by default until auth success.
    const statusBox = document.getElementById('users-status-box');
    if (statusBox) {
        statusBox.classList.add('hidden');
        statusBox.style.display = 'none';
        statusBox.style.visibility = 'hidden';
    }

    // DOM contract: chat interface starts hidden; users list visible for selection.
    const chatInterface = document.getElementById('chat-interface');
    const usersList = document.getElementById('users-status-list');

    if (chatInterface) {
        chatInterface.classList.add('hidden');
        chatInterface.style.display = 'none';
    }

    if (usersList) {
        usersList.style.display = 'block';
        usersList.style.visibility = 'visible';
    }

    // Order matters: slight delay for cookies post-redirect; auth before router for correct page; WebSocket early for realtime handlers.
    setTimeout(() => {
        checkAuthStatus();
        setupWebSocket();
        setupRouter();
        setupEventListeners();
        setupUsersStatusBoxDrag();
    }, 100);
});

// ============================================================================
// SECTION: AUTH VISIBILITY HANDLER
// ============================================================================
// Check auth when page becomes visible (handles login redirects)
document.addEventListener('visibilitychange', () => {
    if (!document.hidden) {
        checkAuthStatus();
    }
});

// ============================================================================
// SECTION: AUTHENTICATION FLOW WIRING
// ============================================================================
// Check authentication status
async function checkAuthStatus() {
    try {
        const response = await fetch('/api/auth/status', {
            credentials: 'include' // Include cookies
        });
        const data = await response.json();
        console.log('Auth status:', data);
        console.log('Logged in?', data.loggedIn);
        if (data.loggedIn) {
            AppState.currentUser = {
                id: data.userID,
                username: data.username
            };
            console.log('User logged in:', AppState.currentUser);
            console.log('Showing users status box and updating navigation...');
            showUsersStatusBox();
            loadUsersStatus();
            loadNotifications(); // Load notifications when user logs in
        } else {
            console.log('User not logged in');
            AppState.currentUser = null;
            hideUsersStatusBox();
            AppState.unreadNotifications = 0;
            AppState.notifications = [];
        }
        updateNavigation();
        // Always re-render to update UI with new auth state
        renderPage();
    } catch (error) {
        console.error('Auth check failed:', error);
    }
}

// ============================================================================
// SECTION: WEBSOCKET LIFECYCLE & ROUTING
// ============================================================================
// Setup WebSocket connection
function setupWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    AppState.ws = new WebSocket(wsUrl);

    AppState.ws.onopen = () => {
        console.log('WebSocket connected');
    };

    AppState.ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            handleWebSocketMessage(data);
        } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
        }
    };

    AppState.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };

    AppState.ws.onclose = () => {
        console.log('WebSocket disconnected, reconnecting...');
        setTimeout(setupWebSocket, 3000);
    };
}

// Handle WebSocket messages
function handleWebSocketMessage(data) {
    console.log('WebSocket message received:', data);

    switch (data.type) {
        case 'new_post':
            if (AppState.currentPage === 'home') {
                loadPosts();
            }
            break;
        case 'new_comment':
            if (AppState.currentPage === 'post' && data.post_id) {
                loadComments(data.post_id);
            }
            break;
        case 'private_message':
            handleNewPrivateMessage(data);
            break;
        case 'user_status_update':
            // Update users status when online status changes
            console.log('User status update received:', data);
            if (AppState.currentUser) {
                loadUsersStatus();
            }
            break;
        case 'notification':
            handleNewNotification(data);
            break;
        case 'typing_start':
            handleTypingStart(data);
            break;
        case 'typing_stop':
            handleTypingStop(data);
            break;
        default:
            console.log('Unknown WebSocket message type:', data.type);
    }
}

// ============================================================================
// SECTION: MESSAGE HANDLING / NOTIFICATIONS
// ============================================================================
// Handle new private message
function handleNewPrivateMessage(data) {
    console.log('Received private message via WebSocket:', data);

    // Get current user ID
    const currentUserId = AppState.currentUser ? AppState.currentUser.id : null;
    if (!currentUserId) {
        console.log('No current user, ignoring message');
        return;
    }

    // Convert IDs to numbers for comparison
    const receiverId = Number(data.receiver_id);
    const senderId = Number(data.sender_id);
    const currentId = Number(currentUserId);

    // Check if this message is for the current user or from the current user
    const isForMe = receiverId === currentId;
    const isFromMe = senderId === currentId;

    if (!isForMe && !isFromMe) {
        console.log('Message is not for current user, ignoring');
        return;
    }

    // If this is a message FOR the current user (not from them)
    if (isForMe && !isFromMe) {
        const senderUserId = senderId;

        // If we're NOT in a chat with this user, show notification
        if (!AppState.currentChatUser || Number(AppState.currentChatUser.id) !== senderUserId) {
            // Increment unread count
            if (!AppState.unreadMessages[senderUserId]) {
                AppState.unreadMessages[senderUserId] = 0;
            }
            AppState.unreadMessages[senderUserId]++;

            // Show notification
            showMessageNotification(data.sender_username || data.sender_name || 'Someone', data.content, senderUserId);
        }
    }

    // If we're in a chat with this user, reload messages to show the new one
    if (AppState.currentChatUser) {
        const chatUserId = Number(AppState.currentChatUser.id);
        const otherUserId = isForMe ? senderId : receiverId;

        if (chatUserId === otherUserId) {
            // We're chatting with this user, reload messages
            console.log('Reloading messages for active chat');
            AppState.messagesOffset = 0;
            loadMessages(chatUserId, false);
            // Clear unread count for this user since we're viewing the chat
            AppState.unreadMessages[chatUserId] = 0;
        }
    }

    // Always refresh users status to update last message and ordering
    if (AppState.currentUser) {
        loadUsersStatus();
    }
}

// Show message notification
function showMessageNotification(senderName, messageContent, senderId) {
    // Create notification element
    const notification = document.createElement('div');
    notification.style.cssText = `
        position: fixed;
        top: 100px;
        right: 380px;
        width: 320px;
        padding: 16px;
        background: linear-gradient(135deg, rgba(30, 41, 59, 0.95) 0%, rgba(15, 23, 42, 0.95) 100%);
        border: 1px solid rgba(34, 197, 94, 0.3);
        border-radius: 16px;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 0 1px rgba(34, 197, 94, 0.2);
        z-index: 100;
        animation: slideInRight 0.3s ease-out;
        cursor: pointer;
        backdrop-filter: blur(10px);
    `;

    const truncatedContent = messageContent.length > 50 ? messageContent.substring(0, 50) + '...' : messageContent;

    notification.innerHTML = `
        <div style="display: flex; align-items: flex-start; gap: 12px;">
            <div style="width: 40px; height: 40px; border-radius: 50%; background: linear-gradient(135deg, #22c55e 0%, #10b981 100%); display: flex; align-items: center; justify-center; flex-shrink: 0; box-shadow: 0 4px 12px rgba(34, 197, 94, 0.3);">
                <svg style="width: 20px; height: 20px; color: white;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                </svg>
            </div>
            <div style="flex: 1; min-width: 0;">
                <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px;">
                    <h4 style="font-weight: 600; color: #f1f5f9; font-size: 14px; margin: 0;">${escapeHtml(senderName)}</h4>
                    <button class="notification-close" style="background: transparent; border: none; color: #94a3b8; cursor: pointer; padding: 4px; font-size: 18px; line-height: 1;">×</button>
                </div>
                <p style="color: #cbd5e1; font-size: 13px; margin: 0; line-height: 1.4; word-wrap: break-word;">${escapeHtml(truncatedContent)}</p>
            </div>
        </div>
    `;

    // Add click handler to open chat
    notification.onclick = (e) => {
        if (!e.target.classList.contains('notification-close')) {
            // Open chat directly
            openChat(senderId, senderName, true); // Assume online for notification
            removeNotification(notification);
        }
    };

    // Close button handler
    const closeBtn = notification.querySelector('.notification-close');
    if (closeBtn) {
        closeBtn.onclick = (e) => {
            e.stopPropagation();
            removeNotification(notification);
        };
    }

    document.body.appendChild(notification);

    // Auto-remove after 5 seconds
    setTimeout(() => {
        removeNotification(notification);
    }, 5000);
}

// Remove notification
function removeNotification(notification) {
    if (notification && notification.parentNode) {
        notification.style.animation = 'slideOutRight 0.3s ease-out';
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 300);
    }
}

// Handle new notification from WebSocket
function handleNewNotification(data) {
    console.log('Received notification via WebSocket:', data);

    // Increment unread count
    AppState.unreadNotifications++;

    // Add to notifications array
    const notification = {
        id: data.notification_id,
        type: data.notificationType,
        title: data.title,
        message: data.message,
        relatedPostId: data.related_post_id,
        relatedCommentId: data.related_comment_id,
        relatedUserId: data.related_user_id,
        isRead: false,
        timeAgo: 'Just now'
    };

    AppState.notifications.unshift(notification);

    // Update navigation badge
    updateNavigation();

    // Show notification popup
    showNotificationPopup(notification);
}

// Show notification popup
function showNotificationPopup(notification) {
    // Determine icon and background color based on notification type
    let icon, bgGradient, borderColor, shadowColor;
    if (notification.type === 'like') {
        icon = '👍';
        bgGradient = 'linear-gradient(135deg, #22c55e 0%, #10b981 100%)';
        borderColor = 'rgba(34, 197, 94, 0.3)';
        shadowColor = 'rgba(34, 197, 94, 0.3)';
    } else if (notification.type === 'dislike') {
        icon = '👎';
        bgGradient = '#ef4444';
        borderColor = 'rgba(239, 68, 68, 0.3)';
        shadowColor = 'rgba(239, 68, 68, 0.3)';
    } else if (notification.type === 'comment') {
        icon = '💬';
        bgGradient = 'linear-gradient(135deg, #22c55e 0%, #10b981 100%)';
        borderColor = 'rgba(34, 197, 94, 0.3)';
        shadowColor = 'rgba(34, 197, 94, 0.3)';
    } else {
        icon = '🔔';
        bgGradient = 'linear-gradient(135deg, #22c55e 0%, #10b981 100%)';
        borderColor = 'rgba(34, 197, 94, 0.3)';
        shadowColor = 'rgba(34, 197, 94, 0.3)';
    }

    // Create notification element
    const notificationEl = document.createElement('div');
    notificationEl.style.cssText = `
        position: fixed;
        top: 100px;
        right: 380px;
        width: 320px;
        padding: 16px;
        background: linear-gradient(135deg, rgba(30, 41, 59, 0.95) 0%, rgba(15, 23, 42, 0.95) 100%);
        border: 1px solid ${borderColor};
        border-radius: 16px;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 0 0 1px ${borderColor.replace('0.3', '0.2')};
        z-index: 100;
        animation: slideInRight 0.3s ease-out;
        cursor: pointer;
        backdrop-filter: blur(10px);
    `;

    notificationEl.innerHTML = `
        <div style="display: flex; align-items: flex-start; gap: 12px;">
            <div style="width: 40px; height: 40px; border-radius: 50%; background: ${bgGradient}; display: flex; align-items: center; justify-content: center; flex-shrink: 0; box-shadow: 0 4px 12px ${shadowColor}; position: relative;">
                <span style="font-size: 20px; line-height: 1; display: flex; align-items: center; justify-content: center; width: 100%; height: 100%; margin: 0; padding: 0;">${icon}</span>
            </div>
            <div style="flex: 1; min-width: 0;">
                <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px;">
                    <h4 style="font-weight: 600; color: #f1f5f9; font-size: 14px; margin: 0;">${escapeHtml(notification.title)}</h4>
                    <button class="notification-close" style="background: transparent; border: none; color: #94a3b8; cursor: pointer; padding: 4px; font-size: 18px; line-height: 1;">×</button>
                </div>
                <p style="color: #cbd5e1; font-size: 13px; margin: 0; line-height: 1.4; word-wrap: break-word;">${escapeHtml(notification.message)}</p>
            </div>
        </div>
    `;

    // Add click handler to navigate to post
    notificationEl.onclick = (e) => {
        if (!e.target.classList.contains('notification-close')) {
            if (notification.relatedPostId) {
                navigateTo(`/post/${notification.relatedPostId}`);
            }
            removeNotificationPopup(notificationEl);
        }
    };

    // Close button handler
    const closeBtn = notificationEl.querySelector('.notification-close');
    if (closeBtn) {
        closeBtn.onclick = (e) => {
            e.stopPropagation();
            removeNotificationPopup(notificationEl);
        };
    }

    document.body.appendChild(notificationEl);

    // Auto-remove after 5 seconds
    setTimeout(() => {
        removeNotificationPopup(notificationEl);
    }, 5000);
}

// Remove notification popup
function removeNotificationPopup(notification) {
    if (notification && notification.parentNode) {
        notification.style.animation = 'slideOutRight 0.3s ease-out';
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 300);
    }
}

// Load notifications from API
async function loadNotifications() {
    if (!AppState.currentUser) {
        return;
    }

    try {
        const response = await fetch('/api/notifications', {
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Failed to fetch notifications');
        }

        const data = await response.json();

        // Count unread notifications
        AppState.unreadNotifications = data.unread ? data.unread.length : 0;

        // Store all notifications
        AppState.notifications = [
            ...(data.unread || []),
            ...(data.read || [])
        ];

        // Update navigation badge
        updateNavigation();
    } catch (error) {
        console.error('Failed to load notifications:', error);
    }
}

// Show notifications dropdown
function showNotificationsDropdown() {
    // Remove existing dropdown if any
    const existingDropdown = document.getElementById('notifications-dropdown');
    if (existingDropdown) {
        existingDropdown.remove();
        return;
    }

    // Create dropdown
    const dropdown = document.createElement('div');
    dropdown.id = 'notifications-dropdown';
    dropdown.style.cssText = `
        position: fixed;
        top: 80px;
        right: 24px;
        width: 400px;
        max-height: 600px;
        background: linear-gradient(135deg, rgba(15, 23, 42, 0.98) 0%, rgba(30, 41, 59, 0.98) 100%);
        border: 1px solid rgba(148, 163, 184, 0.2);
        border-radius: 16px;
        box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
        z-index: 1000;
        overflow: hidden;
        backdrop-filter: blur(10px);
    `;

    const unreadNotifications = AppState.notifications.filter(n => !n.isRead);
    const readNotifications = AppState.notifications.filter(n => n.isRead);

    dropdown.innerHTML = `
        <div style="padding: 16px; border-bottom: 1px solid rgba(148, 163, 184, 0.1); display: flex; justify-content: space-between; align-items: center;">
            <h3 style="font-weight: 700; color: #f1f5f9; font-size: 18px; margin: 0;">Notifications</h3>
            ${unreadNotifications.length > 0 ? `
                <button onclick="markAllNotificationsRead()" style="background: transparent; border: none; color: #22c55e; cursor: pointer; font-size: 12px; font-weight: 600; padding: 4px 8px; border-radius: 8px; transition: background 0.2s;" onmouseover="this.style.background='rgba(34, 197, 94, 0.1)';" onmouseout="this.style.background='transparent';">
                    Mark all read
                </button>
            ` : ''}
            <button onclick="document.getElementById('notifications-dropdown').remove();" style="background: transparent; border: none; color: #94a3b8; cursor: pointer; padding: 4px; font-size: 20px; line-height: 1;">×</button>
        </div>
        <div style="max-height: 500px; overflow-y: auto;">
            ${unreadNotifications.length === 0 && readNotifications.length === 0 ? `
                <div style="padding: 40px; text-align: center; color: #94a3b8;">
                    <p style="margin: 0;">No notifications yet</p>
                </div>
            ` : ''}
            ${unreadNotifications.length > 0 ? `
                <div style="padding: 8px 16px; background: rgba(34, 197, 94, 0.1); border-bottom: 1px solid rgba(148, 163, 184, 0.1);">
                    <span style="font-size: 11px; font-weight: 700; color: #22c55e; text-transform: uppercase; letter-spacing: 0.5px;">Unread (${unreadNotifications.length})</span>
                </div>
                ${unreadNotifications.map(n => renderNotificationItem(n)).join('')}
            ` : ''}
            ${readNotifications.length > 0 ? `
                <div style="padding: 8px 16px; background: rgba(100, 116, 139, 0.1); border-top: 1px solid rgba(148, 163, 184, 0.1); border-bottom: 1px solid rgba(148, 163, 184, 0.1);">
                    <span style="font-size: 11px; font-weight: 700; color: #94a3b8; text-transform: uppercase; letter-spacing: 0.5px;">Read (${readNotifications.length})</span>
                </div>
                ${readNotifications.map(n => renderNotificationItem(n)).join('')}
            ` : ''}
        </div>
    `;

    document.body.appendChild(dropdown);

    // Close on outside click
    setTimeout(() => {
        document.addEventListener('click', function closeDropdown(e) {
            if (!dropdown.contains(e.target) && !e.target.closest('button[onclick*="showNotificationsDropdown"]')) {
                dropdown.remove();
                document.removeEventListener('click', closeDropdown);
            }
        });
    }, 100);
}

// Render notification item
function renderNotificationItem(notification) {
    const icon = notification.type === 'like' ? '👍' : notification.type === 'comment' ? '💬' : '🔔';
    const bgColor = notification.isRead ? 'rgba(30, 41, 59, 0.3)' : 'rgba(34, 197, 94, 0.1)';
    const borderColor = notification.isRead ? 'rgba(148, 163, 184, 0.1)' : 'rgba(34, 197, 94, 0.3)';

    return `
        <div onclick="handleNotificationClick(${notification.id}, ${notification.relatedPostId || 'null'})" style="padding: 12px 16px; border-bottom: 1px solid rgba(148, 163, 184, 0.1); cursor: pointer; transition: background 0.2s; background: ${bgColor}; border-left: 3px solid ${borderColor};" onmouseover="this.style.background='rgba(30, 41, 59, 0.6)';" onmouseout="this.style.background='${bgColor}';">
            <div style="display: flex; align-items: flex-start; gap: 12px;">
                <div style="font-size: 24px; flex-shrink: 0;">${icon}</div>
                <div style="flex: 1; min-width: 0;">
                    <h4 style="font-weight: 600; color: #f1f5f9; font-size: 14px; margin: 0 0 4px 0;">${escapeHtml(notification.title)}</h4>
                    <p style="color: #cbd5e1; font-size: 13px; margin: 0 0 4px 0; line-height: 1.4;">${escapeHtml(notification.message)}</p>
                    <span style="color: #94a3b8; font-size: 11px;">${notification.timeAgo || 'Just now'}</span>
                </div>
            </div>
        </div>
    `;
}

// Handle notification click
async function handleNotificationClick(notificationId, postId) {
    // Mark as read
    if (!AppState.notifications.find(n => n.id === notificationId)?.isRead) {
        try {
            const formData = new FormData();
            formData.append('notification_id', notificationId);

            await fetch('/api/notifications/mark-read', {
                method: 'POST',
                body: formData,
                credentials: 'include'
            });

            // Update local state
            const notification = AppState.notifications.find(n => n.id === notificationId);
            if (notification) {
                notification.isRead = true;
                AppState.unreadNotifications = Math.max(0, AppState.unreadNotifications - 1);
            }
        } catch (error) {
            console.error('Failed to mark notification as read:', error);
        }
    }

    // Navigate to post if available
    if (postId) {
        navigateTo(`/post/${postId}`);
    }

    // Close dropdown
    const dropdown = document.getElementById('notifications-dropdown');
    if (dropdown) {
        dropdown.remove();
    }

    // Update navigation
    updateNavigation();
}

// Mark all notifications as read
async function markAllNotificationsRead() {
    try {
        await fetch('/api/notifications/mark-all-read', {
            method: 'POST',
            credentials: 'include'
        });

        // Update local state
        AppState.notifications.forEach(n => n.isRead = true);
        AppState.unreadNotifications = 0;

        // Reload notifications
        await loadNotifications();

        // Update navigation
        updateNavigation();

        // Refresh dropdown
        const dropdown = document.getElementById('notifications-dropdown');
        if (dropdown) {
            showNotificationsDropdown();
        }
    } catch (error) {
        console.error('Failed to mark all notifications as read:', error);
    }
}

// Make functions available globally
window.showNotificationsDropdown = showNotificationsDropdown;
window.handleNotificationClick = handleNotificationClick;
window.markAllNotificationsRead = markAllNotificationsRead;

// Router setup
function setupRouter() {
    window.addEventListener('popstate', () => {
        const path = window.location.pathname;
        navigateTo(path);
    });

    // Initial navigation
    navigateTo(window.location.pathname);
}

// Navigation function (for initial load and popstate)
function navigateTo(path) {
    navigateInternal(path);
}

// Update navigation bar
function updateNavigation() {
    const navLinks = document.getElementById('nav-links');
    if (!navLinks) {
        console.warn('Navigation links element not found');
        return;
    }

    navLinks.innerHTML = '';

    if (AppState.currentUser) {
        console.log('Updating navigation for logged in user:', AppState.currentUser.username);
        const notificationBadge = AppState.unreadNotifications > 0
            ? `<span class="notification-badge" style="position: absolute; top: -6px; right: -6px; background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); color: white; border-radius: 10px; min-width: 20px; height: 20px; padding: 0 6px; display: flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 700; box-shadow: 0 2px 8px rgba(239, 68, 68, 0.4);">${AppState.unreadNotifications > 99 ? '99+' : AppState.unreadNotifications}</span>`
            : '';
        navLinks.innerHTML = `
            <a href="#" onclick="window.navigateTo('/create-post'); return false;" class="btn-primary">Create Post</a>
            <div style="position: relative; display: inline-block;">
                <button onclick="showNotificationsDropdown()" style="background: transparent; border: none; cursor: pointer; padding: 8px; position: relative; color: #64748b; transition: color 0.2s;" onmouseover="this.style.color='#f1f5f9';" onmouseout="this.style.color='#64748b';" title="Notifications">
                    <svg style="width: 20px; height: 20px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"></path>
                    </svg>
                    ${notificationBadge}
                </button>
            </div>
            <span class="text-gray-700 font-semibold" style="display: inline-block;">${escapeHtml(AppState.currentUser.username)}</span>
            <a href="/logout" class="btn-secondary">Logout</a>
        `;
    } else {
        navLinks.innerHTML = `
            <a href="#" onclick="window.navigateTo('/login'); return false;" class="btn-primary">Login</a>
            <a href="#" onclick="window.navigateTo('/register'); return false;" class="btn-secondary">Register</a>
        `;
    }
}

// Render current page
function renderPage() {
    const app = document.getElementById('app');

    // Check if user needs to be logged in
    if (['create-post', 'post'].includes(AppState.currentPage) && !AppState.currentUser) {
        navigateTo('/login');
        return;
    }

    switch (AppState.currentPage) {
        case 'home':
            renderHome();
            break;
        case 'login':
            renderLogin();
            break;
        case 'register':
            renderRegister();
            break;
        case 'post':
            renderPost();
            break;
        case 'create-post':
            renderCreatePost();
            break;
        case 'profile':
            renderProfile();
            break;
        default:
            renderHome();
    }

    updateNavigation();
}

// Render home page
async function renderHome() {
    const app = document.getElementById('app');

    // Load categories for filter
    let categories = [];
    try {
        const response = await fetch('/api/categories', {
            credentials: 'include'
        });
        if (response.ok) {
            categories = await response.json();
        }
    } catch (error) {
        console.error('Failed to load categories:', error);
    }

    app.innerHTML = `
        <div class="mb-10">
            <div class="mb-6">
                <h1 class="text-5xl font-extrabold gradient-text mb-3 tracking-tight">Welcome to Tech Talk Forum</h1>
                <p class="text-gray-300 text-lg">Share your thoughts and connect with others! </p>
            </div>
            
            <!-- Categories Filter -->
            ${categories.length > 0 ? `
                <div class="mb-6">
                    <h3 class="text-sm font-semibold text-gray-400 uppercase tracking-wide mb-3">Filter by Category</h3>
                    <div class="flex flex-wrap gap-2" id="category-filters">
                        <button onclick="loadPosts()" class="category-filter-btn active" data-category="all">
                            All Posts
                        </button>
                        ${categories.map(cat => `
                            <button onclick="loadPostsByCategory('${escapeHtml(cat.name || cat.category_name || cat)}')" class="category-filter-btn" data-category="${escapeHtml(cat.name || cat.category_name || cat)}">
                                ${escapeHtml(cat.name || cat.category_name || cat)}
                            </button>
                        `).join('')}
                    </div>
                </div>
            ` : ''}
        </div>
        <div id="posts-container" class="space-y-4">
            <div class="text-center py-16">
                <div class="animate-spin rounded-full h-12 w-12 border-4 border-slate-300 border-t-slate-500 mx-auto mb-4"></div>
                <p class="text-gray-500 font-medium">Loading posts...</p>
            </div>
        </div>
    `;
    await loadPosts();
}

// Load posts by category
async function loadPostsByCategory(categoryName) {
    const container = document.getElementById('posts-container');
    if (!container) return;

    // Update active filter button
    document.querySelectorAll('.category-filter-btn').forEach(btn => {
        btn.classList.remove('active');
        if (btn.dataset.category === categoryName) {
            btn.classList.add('active');
        }
    });

    try {
        const response = await fetch(`/api/posts?filter=categories&value=${encodeURIComponent(categoryName)}`, {
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const posts = await response.json();

        if (posts === null || posts === undefined) {
            renderPosts([]);
            return;
        }

        if (!Array.isArray(posts)) {
            throw new Error('Invalid data format received');
        }

        renderPosts(posts);
    } catch (error) {
        console.error('Failed to load posts:', error);
        container.innerHTML = `
            <div class="card text-center py-12">
                <div class="text-6xl mb-4">😕</div>
                <p class="text-red-500 text-lg font-semibold mb-2">Failed to load posts</p>
                <p class="text-sm text-gray-500 mb-4">${error.message || 'Please try again later'}</p>
                <button onclick="loadPosts()" class="btn-primary">🔄 Retry</button>
            </div>
        `;
    }
}

// Load posts
async function loadPosts() {
    // Update active filter button
    document.querySelectorAll('.category-filter-btn')?.forEach(btn => {
        btn.classList.remove('active');
        if (btn.dataset.category === 'all') {
            btn.classList.add('active');
        }
    });
    const container = document.getElementById('posts-container');
    if (!container) return;

    try {
        const response = await fetch('/api/posts', {
            credentials: 'include'
        });

        // Log response for debugging
        console.log('Response status:', response.status);
        console.log('Response headers:', response.headers.get('content-type'));

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Error response:', errorText);
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const contentType = response.headers.get('content-type');
        if (!contentType || !contentType.includes('application/json')) {
            const text = await response.text();
            console.error('Non-JSON response:', text);
            throw new Error('Invalid response format');
        }

        const posts = await response.json();
        console.log('Received posts:', posts);

        // Handle null response
        if (posts === null || posts === undefined) {
            console.warn('Posts is null/undefined, treating as empty array');
            renderPosts([]);
            return;
        }

        // Ensure posts is an array
        if (!Array.isArray(posts)) {
            console.error('Posts is not an array:', posts, typeof posts);
            throw new Error('Invalid data format received');
        }

        renderPosts(posts);
    } catch (error) {
        console.error('Failed to load posts:', error);
        container.innerHTML = `
            <div class="card text-center py-12">
                <div class="text-6xl mb-4">😕</div>
                <p class="text-red-500 text-lg font-semibold mb-2">Failed to load posts</p>
                <p class="text-sm text-gray-500 mb-4">${error.message || 'Please try again later'}</p>
                <button onclick="loadPosts()" class="btn-primary">🔄 Retry</button>
            </div>
        `;
    }
}

// Render posts
function renderPosts(posts) {
    const container = document.getElementById('posts-container');
    if (!container) return;

    // Handle null/undefined posts
    if (!posts || !Array.isArray(posts)) {
        container.innerHTML = `
            <div class="card text-center py-12">
                <p class="text-red-500">Error: Invalid data received</p>
                <button onclick="loadPosts()" class="btn-primary mt-4">Retry</button>
            </div>
        `;
        return;
    }

    if (posts.length === 0) {
        container.innerHTML = `
            <div class="card text-center py-16">
                <div class="text-7xl mb-6 animate-bounce">📝</div>
                <h3 class="text-2xl font-bold text-gray-100 mb-2">No posts yet</h3>
                <p class="text-gray-400 mb-6">Be the first to share your thoughts!</p>
                ${AppState.currentUser ? `<a href="#" onclick="navigateTo('/create-post'); return false;" class="btn-primary inline-flex items-center gap-2">
                    <span></span> Create First Post
                </a>` : ''}
            </div>
        `;
        return;
    }

    container.innerHTML = posts.map((post, index) => `
        <article class="post-card group" style="animation: fadeInUp 0.5s ease-out ${index * 0.06}s both;" onclick="navigateTo('/post/${post.id}')">
            <div class="post-card-content">
                <div class="flex items-center justify-between mb-4">
                    <div class="flex items-center gap-3">
                        <div class="post-author-avatar flex-shrink-0">
                            <div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 via-slate-400 to-slate-500 flex items-center justify-center text-white shadow-xl ring-2 ring-slate-400/40 hover:ring-slate-400/60 transition-all">
                                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                                </svg>
                            </div>
                        </div>
                        <span class="font-bold text-gray-50 text-base leading-tight">${escapeHtml(post.author || 'Unknown')}</span>
                    </div>
                    <div class="flex items-center gap-1.5 text-gray-400 flex-shrink-0" style="font-size: 10px; white-space: nowrap;">
                        <svg class="w-3 h-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                        </svg>
                        <span class="text-gray-400" style="font-size: 10px; font-weight: 400; white-space: nowrap;">${post.timeAgo || 'Just now'}</span>
                    </div>
                </div>
                <div class="flex items-start gap-5 mb-5">
                    <div class="flex-1 min-w-0">
                        <h2 class="text-2xl font-bold text-gray-50 mb-3 group-hover:text-slate-300 transition-colors leading-tight line-clamp-2">${escapeHtml(post.title || 'Untitled')}</h2>
                        <p class="text-gray-300 mb-4 line-clamp-3 leading-relaxed text-base">${escapeHtml((post.excerpt || post.content || '').substring(0, 250))}${(post.content && post.content.length > 250) ? '...' : ''}</p>
                    </div>
                    ${post.thumbnailUrl ? `
                        <div class="post-thumbnail">
                            <img src="${post.thumbnailUrl}" alt="Post image" class="w-32 h-32 object-cover rounded-2xl shadow-xl ring-2 ring-slate-400/20 group-hover:ring-slate-400/40 transition-all group-hover:scale-105">
                        </div>
                    ` : ''}
                </div>
                
                ${(post.tags && Array.isArray(post.tags) && post.tags.length > 0) ? `
                    <div class="mb-5 flex flex-wrap gap-2">
                        ${post.tags.map(tag => `
                            <span class="post-tag">
                                <span class="tag-icon">#</span>
                                ${escapeHtml(tag)}
                            </span>
                        `).join('')}
                    </div>
                ` : ''}
                
                <div class="post-stats">
                    <button class="stat-item stat-likes ${post.userVote === 1 ? 'active' : ''}" title="Like" onclick="event.stopPropagation(); handlePostVote(${post.id}, 1, 'post-card')">
                        <span class="stat-emoji">👍</span>
                        <div class="stat-icon">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5"></path>
                            </svg>
                        </div>
                        <span class="stat-value" id="post-likes-${post.id}">${post.likes || 0}</span>
                    </button>
                    <button class="stat-item stat-dislikes ${post.userVote === -1 ? 'active' : ''}" title="Dislike" onclick="event.stopPropagation(); handlePostVote(${post.id}, -1, 'post-card')">
                        <span class="stat-emoji">👎</span>
                        <div class="stat-icon">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14H5.236a2 2 0 01-1.789-2.894l3.5-7A2 2 0 018.736 3h4.018a2 2 0 01.485.06l3.76.94m-7 10v5a2 2 0 002 2h.096c.5 0 .905-.405.905-.904 0-.715.211-1.413.608-2.008L17 13V4m-7 10h2m5-10h2a2 2 0 012 2v6a2 2 0 01-2 2h-2.5"></path>
                            </svg>
                        </div>
                        <span class="stat-value" id="post-dislikes-${post.id}">${post.dislikes || 0}</span>
                    </button>
                    <div class="stat-item stat-comments" title="Comments">
                        <span class="stat-emoji">💬</span>
                        <div class="stat-icon">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
                            </svg>
                        </div>
                        <span class="stat-value">${post.comments || 0}</span>
                    </div>
                    <div class="stat-item stat-read-more">
                        <span class="read-more-text">Read more</span>
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                        </svg>
                    </div>
                </div>
            </div>
        </article>
    `).join('');
}

// Render login page
function renderLogin() {
    const app = document.getElementById('app');

    // Check for error in URL
    const urlParams = new URLSearchParams(window.location.search);
    const error = urlParams.get('error');
    const errorMessage = error ? decodeURIComponent(error.replace(/\+/g, ' ')) : '';

    app.innerHTML = `
        <div class="max-w-md mx-auto mt-12">
            <div class="card">
                <h2 class="text-3xl font-bold gradient-text mb-6 text-center">Login</h2>
                
                <form id="login-form" method="POST" action="/login" class="space-y-4">
                    <div>
                        <label class="block text-gray-700 font-semibold mb-2">Email or Username</label>
                        <input type="text" name="email" required class="input-field" placeholder="Enter your email or username">
                    </div>
                    <div>
                        <label class="block text-gray-700 font-semibold mb-2">Password</label>
                        <input type="password" name="password" required class="input-field" placeholder="Enter your password">
                    </div>
                    ${errorMessage ? `<div id="login-error" class="text-red-500 text-sm">${escapeHtml(errorMessage)}</div>` : '<div id="login-error" class="text-red-500 text-sm hidden"></div>'}
                    <button type="submit" class="btn-primary w-full">Login</button>
                </form>
                <p class="mt-4 text-center text-gray-600">
                    Don't have an account? <a href="#" onclick="navigateTo('/register'); return false;" class="text-slate-400 font-semibold hover:underline">Register</a>
                </p>
            </div>
        </div>
    `;

    // Clean URL if there was an error
    if (errorMessage) {
        window.history.replaceState({}, '', '/login');
    }

    // Attach event listener to login form
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
}

// Handle login
async function handleLogin(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const errorDiv = document.getElementById('login-error');

    // Clear previous errors
    errorDiv.classList.add('hidden');
    errorDiv.textContent = '';

    try {
        const response = await fetch('/login', {
            method: 'POST',
            body: formData,
            credentials: 'include',
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        });

        // Login redirects on success (303), so if we get here it's an error
        if (response.status >= 200 && response.status < 300) {
            // Success - redirect happened, reload page
            window.location.href = '/';
            return;
        }

        // Handle error response
        if (response.headers.get('content-type')?.includes('application/json')) {
            const data = await response.json();
            errorDiv.textContent = data.error || 'Login failed. Please try again.';
        } else {
            // HTML response (shouldn't happen with AJAX header, but handle it)
            errorDiv.textContent = 'Invalid credentials. Please try again.';
        }
        errorDiv.classList.remove('hidden');
    } catch (error) {
        console.error('Login error:', error);
        errorDiv.textContent = 'Login failed. Please try again.';
        errorDiv.classList.remove('hidden');
    }
}

// Render register page
function renderRegister() {
    const app = document.getElementById('app');
    app.innerHTML = `
        <div class="max-w-md mx-auto mt-8">
            <div class="card">
                <div class="text-center mb-8">
                    <h2 class="text-3xl font-extrabold gradient-text mb-2">Create Account</h2>
                    <p class="text-gray-500">Join our community today</p>
                </div>
                
                <form id="register-form" class="space-y-5">
                    <div>
                        <label class="block text-sm font-semibold text-gray-700 mb-2">Nickname</label>
                        <input type="text" name="username" required class="input-field" placeholder="Choose a nickname">
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-sm font-semibold text-gray-700 mb-2">First Name</label>
                            <input type="text" name="first_name" required class="input-field" placeholder="First name">
                        </div>
                        <div>
                            <label class="block text-sm font-semibold text-gray-700 mb-2">Last Name</label>
                            <input type="text" name="last_name" required class="input-field" placeholder="Last name">
                        </div>
                    </div>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="block text-sm font-semibold text-gray-700 mb-2">Age</label>
                            <input type="number" name="age" required min="13" class="input-field" placeholder="Age">
                        </div>
                        <div>
                            <label class="block text-sm font-semibold text-gray-700 mb-2">Gender</label>
                            <select name="gender" required class="input-field">
                                <option value="">Select</option>
                                <option value="Male">Male</option>
                                <option value="Female">Female</option>
                                <option value="Other">Other</option>
                            </select>
                        </div>
                    </div>
                    <div>
                        <label class="block text-sm font-semibold text-gray-700 mb-2">Email</label>
                        <input type="email" name="email" required class="input-field" placeholder="Enter your email">
                    </div>
                    <div>
                        <label class="block text-sm font-semibold text-gray-700 mb-2">Password</label>
                        <input type="password" name="password" required class="input-field" placeholder="Create a password">
                        <p class="text-xs text-gray-400 mt-1.5">Must have 8+ chars, uppercase, lowercase, number, and symbol</p>
                    </div>
                    <div>
                        <label class="block text-sm font-semibold text-gray-700 mb-2">Confirm Password</label>
                        <input type="password" name="confirm_password" required class="input-field" placeholder="Confirm password">
                    </div>
                    <div id="register-error" class="text-red-500 text-sm hidden"></div>
                    <button type="submit" class="btn-primary w-full">Create Account</button>
                </form>
                <p class="mt-6 text-center text-sm text-gray-600">
                    Already have an account? <a href="#" onclick="navigateTo('/login'); return false;" class="text-slate-400 font-semibold hover:text-slate-300 hover:underline transition-colors">Sign in</a>
                </p>
            </div>
        </div>
    `;

    document.getElementById('register-form').addEventListener('submit', handleRegister);
}

// Handle register
async function handleRegister(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const errorDiv = document.getElementById('register-error');

    try {
        const response = await fetch('/register', {
            method: 'POST',
            body: formData
        });

        if (response.ok || response.redirected) {
            navigateTo('/login');
        } else {
            const text = await response.text();
            errorDiv.textContent = 'Registration failed. Please check your information.';
            errorDiv.classList.remove('hidden');
        }
    } catch (error) {
        errorDiv.textContent = 'Registration failed. Please try again.';
        errorDiv.classList.remove('hidden');
    }
}

// Render post view page
async function renderPost() {
    const app = document.getElementById('app');
    const postId = AppState.currentPostId;

    app.innerHTML = `
        <div id="post-container" class="mb-8">
            <div class="text-center py-12">
                <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-slate-400 mx-auto"></div>
                <p class="mt-4 text-gray-600">Loading post...</p>
            </div>
        </div>
        <div id="comments-container" class="mt-8">
            <h3 class="text-2xl font-bold mb-4">Comments</h3>
            <div class="text-center py-8">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-slate-400 mx-auto"></div>
            </div>
        </div>
    `;

    await loadPost(postId);
    await loadComments(postId);
}

// Load single post
async function loadPost(postId) {
    try {
        const response = await fetch(`/api/post?id=${postId}`);
        const post = await response.json();
        renderPostDetails(post);
    } catch (error) {
        console.error('Failed to load post:', error);
        document.getElementById('post-container').innerHTML = `
            <div class="card text-center py-12">
                <p class="text-red-500">Failed to load post. Please try again.</p>
            </div>
        `;
    }
}

// Render post details
function renderPostDetails(post) {
    const container = document.getElementById('post-container');
    container.innerHTML = `
        <article class="post-detail-card">
            <div class="post-detail-header">
                <div class="flex items-center justify-between mb-6">
                    <div class="flex items-center gap-3">
                        <div class="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 via-slate-400 to-slate-500 flex items-center justify-center text-white shadow-xl ring-2 ring-slate-400/40 hover:ring-slate-400/60 transition-all flex-shrink-0">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                            </svg>
                        </div>
                        <span class="font-bold text-gray-50 text-lg leading-tight">${escapeHtml(post.author || 'Unknown')}</span>
                    </div>
                    <div class="flex items-center gap-1.5 text-gray-400 flex-shrink-0" style="font-size: 10px; white-space: nowrap;">
                        <svg class="w-3 h-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                        </svg>
                        <span class="text-gray-400" style="font-size: 10px; font-weight: 400; white-space: nowrap;">${post.timeAgo || 'Just now'}</span>
                    </div>
                </div>
                <h1 class="text-4xl font-bold text-gray-50 mb-6 leading-tight">${escapeHtml(post.title)}</h1>
                ${post.imageUrl ? `<img src="${post.imageUrl}" alt="Post image" class="w-full rounded-2xl mb-6 shadow-xl ring-2 ring-slate-400/20">` : ''}
            </div>
            
            <div class="post-detail-content">
                <div class="prose prose-invert max-w-none mb-6">
                    <p class="text-gray-200 whitespace-pre-wrap leading-relaxed text-lg">${escapeHtml(post.content)}</p>
                </div>
                
                ${post.tags && post.tags.length > 0 ? `
                    <div class="flex flex-wrap gap-2 mb-6">
                        ${post.tags.map(tag => `
                            <span class="post-tag">
                                <span class="tag-icon">#</span>
                                ${escapeHtml(tag)}
                            </span>
                        `).join('')}
                    </div>
                ` : ''}
            </div>
            
            <div class="post-detail-stats">
                <button onclick="handlePostVote(${post.id}, 1, 'detail')" class="post-vote-btn post-vote-like ${post.userVote === 1 ? 'active' : ''}" title="Like">
                    <span class="vote-emoji">👍</span>
                    <div class="vote-icon">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5"></path>
                        </svg>
                    </div>
                    <span class="vote-value" id="likes-count-${post.id}">${post.likes || 0}</span>
                </button>
                <button onclick="handlePostVote(${post.id}, -1, 'detail')" class="post-vote-btn post-vote-dislike ${post.userVote === -1 ? 'active' : ''}" title="Dislike">
                    <span class="vote-emoji">👎</span>
                    <div class="vote-icon">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14H5.236a2 2 0 01-1.789-2.894l3.5-7A2 2 0 018.736 3h4.018a2 2 0 01.485.06l3.76.94m-7 10v5a2 2 0 002 2h.096c.5 0 .905-.405.905-.904 0-.715.211-1.413.608-2.008L17 13V4m-7 10h2m5-10h2a2 2 0 012 2v6a2 2 0 01-2 2h-2.5"></path>
                        </svg>
                    </div>
                    <span class="vote-value" id="dislikes-count-${post.id}">${post.dislikes || 0}</span>
                </button>
                <div class="post-comments-count">
                    <span class="comments-emoji">💬</span>
                    <div class="comments-icon">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
                        </svg>
                    </div>
                    <span class="comments-value">${post.comments || 0} comments</span>
                </div>
            </div>
        </article>
    `;
}

// Load comments
async function loadComments(postId) {
    try {
        const response = await fetch(`/api/comments?post_id=${postId}`, {
            credentials: 'include'
        });
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        const comments = await response.json();
        renderComments(comments, postId);
    } catch (error) {
        console.error('Failed to load comments:', error);
        const container = document.getElementById('comments-container');
        if (container) {
            container.innerHTML = `
                <div class="text-center py-8">
                    <p class="text-red-500 mb-4">Failed to load comments</p>
                    <button onclick="loadComments(${postId})" class="btn-primary">Retry</button>
                </div>
            `;
        }
    }
}

// Render comments
function renderComments(comments, postId) {
    const container = document.getElementById('comments-container');
    if (!container) {
        console.error('Comments container not found');
        return;
    }

    // Ensure comments is an array
    if (!Array.isArray(comments)) {
        comments = [];
    }

    // Check if user is logged in
    const isLoggedIn = AppState.currentUser && AppState.currentUser.id;

    container.innerHTML = `
        <div class="flex items-center justify-between mb-6">
            <h3 class="text-2xl font-bold text-gray-100">Comments (${comments.length})</h3>
        </div>
        ${isLoggedIn ? `
            <div class="comment-form-card mb-6">
                <div class="flex items-center gap-3 mb-4">
                    <div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 via-slate-400 to-slate-500 flex items-center justify-center text-white text-sm font-bold shadow-lg ring-2 ring-slate-400/30">
                        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                        </svg>
                    </div>
                    <span class="font-semibold text-gray-200">${escapeHtml(AppState.currentUser.username)}</span>
                </div>
                <form id="comment-form" class="space-y-4">
                    <div>
                        <label class="block text-gray-300 font-semibold mb-2 text-sm">Write a comment</label>
                        <textarea 
                            name="content" 
                            id="comment-content" 
                            required 
                            class="input-field comment-textarea w-full" 
                            rows="5" 
                            placeholder="Share your thoughts..."></textarea>
                    </div>
                    <div class="flex justify-end">
                        <button type="submit" class="btn-primary px-6 py-2.5 flex items-center gap-2">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"></path>
                            </svg>
                            <span>Post Comment</span>
                        </button>
                    </div>
                </form>
            </div>
        ` : `
            <div class="card mb-6 text-center py-6">
                <p class="text-gray-300 mb-3">Please login to comment</p>
                <a href="#" onclick="navigateTo('/login'); return false;" class="btn-primary inline-block">Login</a>
            </div>
        `}
        <div id="comments-list" class="space-y-4">
            ${comments.length === 0 ? '<p class="text-gray-400 text-center py-8">No comments yet. Be the first to comment!</p>' : ''}
        </div>
    `;

    if (comments.length > 0) {
        const commentsList = document.getElementById('comments-list');
        if (commentsList) {
            commentsList.innerHTML = comments.map(comment => `
                <div class="card">
                    <div class="flex items-start justify-between mb-3">
                        <div class="flex-1">
                            <div class="flex items-center gap-2 mb-3">
                                <span class="font-bold text-gray-50 text-base">${escapeHtml(comment.author)}</span>
                            </div>
                            <p class="text-gray-200 whitespace-pre-wrap leading-relaxed">${escapeHtml(comment.content)}</p>
                        </div>
                        <span class="text-gray-400 flex items-center gap-1 flex-shrink-0 ml-3" style="font-size: 10px; white-space: nowrap;">
                            <svg class="w-3 h-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                            </svg>
                            <span style="font-size: 10px; font-weight: 400; white-space: nowrap;">${comment.timeAgo || 'Just now'}</span>
                        </span>
                    </div>
                    <div class="flex items-center space-x-4 mt-4 pt-4 border-t border-gray-200">
                        <button onclick="handleCommentVote(${comment.id}, 1)" class="comment-vote-btn comment-vote-like ${comment.userVote === 1 ? 'active' : ''}">
                            <span class="comment-emoji">👍</span>
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5"></path>
                            </svg>
                            <span class="text-sm font-semibold" id="comment-likes-${comment.id}">${comment.likeCount || 0}</span>
                        </button>
                        <button onclick="handleCommentVote(${comment.id}, -1)" class="comment-vote-btn comment-vote-dislike ${comment.userVote === -1 ? 'active' : ''}">
                            <span class="comment-emoji">👎</span>
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14H5.236a2 2 0 01-1.789-2.894l3.5-7A2 2 0 018.736 3h4.018a2 2 0 01.485.06l3.76.94m-7 10v5a2 2 0 002 2h.096c.5 0 .905-.405.905-.904 0-.715.211-1.413.608-2.008L17 13V4m-7 10h2m5-10h2a2 2 0 012 2v6a2 2 0 01-2 2h-2.5"></path>
                            </svg>
                            <span class="text-sm font-semibold" id="comment-dislikes-${comment.id}">${comment.dislikeCount || 0}</span>
                        </button>
                    </div>
                </div>
            `).join('');
        }
    }

    // Setup comment form event listener
    if (isLoggedIn) {
        // Use requestAnimationFrame to ensure DOM is ready
        requestAnimationFrame(() => {
            const commentForm = document.getElementById('comment-form');
            if (commentForm) {
                // Remove any existing listeners by cloning
                const newForm = commentForm.cloneNode(true);
                commentForm.parentNode.replaceChild(newForm, commentForm);

                // Add event listener to the new form
                const form = document.getElementById('comment-form');
                if (form) {
                    form.addEventListener('submit', (e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        handleCreateComment(e, postId);
                        return false;
                    });
                    console.log('Comment form event listener attached for post', postId);
                } else {
                    console.error('Comment form not found after cloning');
                }
            } else {
                console.warn('Comment form not found - user may not be logged in');
            }
        });
    }
}

// Handle create comment
async function handleCreateComment(e, postId) {
    e.preventDefault();

    if (!AppState.currentUser) {
        console.error('User not logged in');
        return;
    }

    const form = e.target;
    const formData = new FormData(form);
    const content = formData.get('content');

    if (!content || !content.trim()) {
        console.error('Comment content is required');
        return;
    }

    formData.append('post_id', postId);

    // Disable submit button during request
    const submitBtn = form.querySelector('button[type="submit"]');
    const originalText = submitBtn.innerHTML;
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<span class="animate-spin">⏳</span> Posting...';

    try {
        const response = await fetch('/api/comments/create', {
            method: 'POST',
            body: formData,
            credentials: 'include'
        });

        if (response.ok) {
            form.reset();
            await loadComments(postId);
        } else {
            const errorText = await response.text();
            console.error('Failed to create comment:', errorText);
            alert('Failed to post comment. Please try again.');
        }
    } catch (error) {
        console.error('Failed to create comment:', error);
        alert('Failed to post comment. Please try again.');
    } finally {
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalText;
    }
}

// Handle post vote
async function handlePostVote(postId, vote, context = 'detail') {
    if (!AppState.currentUser) {
        navigateTo('/login');
        return;
    }

    const formData = new FormData();
    formData.append('post_id', postId);
    formData.append('vote', vote);

    try {
        const response = await fetch('/api/posts/like', {
            method: 'POST',
            body: formData,
            credentials: 'include'
        });
        const data = await response.json();

        // Update counts based on context
        if (context === 'post-card') {
            // Update post card stats
            const likesEl = document.getElementById(`post-likes-${postId}`);
            const dislikesEl = document.getElementById(`post-dislikes-${postId}`);
            if (likesEl) likesEl.textContent = data.likeCount || 0;
            if (dislikesEl) dislikesEl.textContent = data.dislikeCount || 0;

            // Update active states
            const likeBtn = document.querySelector(`#post-likes-${postId}`)?.closest('.stat-likes');
            const dislikeBtn = document.querySelector(`#post-dislikes-${postId}`)?.closest('.stat-dislikes');

            if (likeBtn && dislikeBtn) {
                // Remove active from both
                likeBtn.classList.remove('active');
                dislikeBtn.classList.remove('active');

                // Add active to the one that was voted
                if (data.userVote === 1) {
                    likeBtn.classList.add('active');
                } else if (data.userVote === -1) {
                    dislikeBtn.classList.add('active');
                }
            }
        } else {
            // Update detail view stats
            const likesEl = document.getElementById(`likes-count-${postId}`);
            const dislikesEl = document.getElementById(`dislikes-count-${postId}`);
            if (likesEl) likesEl.textContent = data.likeCount || 0;
            if (dislikesEl) dislikesEl.textContent = data.dislikeCount || 0;

            // Update active states in detail view
            const likeBtn = document.querySelector(`.post-vote-like[onclick*="${postId}"]`);
            const dislikeBtn = document.querySelector(`.post-vote-dislike[onclick*="${postId}"]`);

            if (likeBtn && dislikeBtn) {
                likeBtn.classList.remove('active');
                dislikeBtn.classList.remove('active');

                if (data.userVote === 1) {
                    likeBtn.classList.add('active');
                } else if (data.userVote === -1) {
                    dislikeBtn.classList.add('active');
                }
            }
        }
    } catch (error) {
        console.error('Failed to vote:', error);
    }
}

// Handle comment vote
async function handleCommentVote(commentId, vote) {
    if (!AppState.currentUser) {
        navigateTo('/login');
        return;
    }

    const formData = new FormData();
    formData.append('comment_id', commentId);
    formData.append('vote', vote);

    try {
        const response = await fetch('/api/comments/like', {
            method: 'POST',
            body: formData
        });
        const data = await response.json();

        document.getElementById(`comment-likes-${commentId}`).textContent = data.likeCount;
        document.getElementById(`comment-dislikes-${commentId}`).textContent = data.dislikeCount;
    } catch (error) {
        console.error('Failed to vote:', error);
    }
}

// Render profile page
async function renderProfile() {
    const app = document.getElementById('app');

    app.innerHTML = `
        <div class="mb-10">
            <div class="text-center py-12">
                <div class="animate-spin rounded-full h-12 w-12 border-4 border-slate-300 border-t-slate-500 mx-auto mb-4"></div>
                <p class="text-gray-500 font-medium">Loading profile...</p>
            </div>
        </div>
    `;

    try {
        const response = await fetch('/api/user/profile', {
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Failed to load profile');
        }

        const profile = await response.json();
        renderProfileDetails(profile);
    } catch (error) {
        console.error('Failed to load profile:', error);
        app.innerHTML = `
            <div class="mb-10">
                <div class="card text-center py-8">
                    <div class="text-2xl mb-3">😕</div>
                    <p class="text-red-500 font-semibold mb-2">Failed to load profile</p>
                    <p class="text-sm text-gray-500 mb-4">${error.message || 'Please try again later'}</p>
                    <button onclick="navigateTo('/profile')" class="btn-primary">🔄 Retry</button>
                </div>
            </div>
        `;
    }
}

// Render profile details with modern layout
function renderProfileDetails(profile) {
    const app = document.getElementById('app');

    const joinDate = profile.joinDate ? new Date(profile.joinDate).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    }) : 'Unknown';

    app.innerHTML = `
        <div class="mb-10">
            <div class="mb-6">
                <h1 class="text-5xl font-extrabold gradient-text mb-3 tracking-tight">User Profile</h1>
                <p class="text-gray-300 text-lg">View your profile information and activity ✨</p>
            </div>
            
            <!-- Profile Header Card -->
            <div class="card mb-6">
                <div class="flex flex-col md:flex-row items-center md:items-start gap-4">
                    <div class="flex-shrink-0">
                        ${profile.profileImage ? `
                            <img src="${profile.profileImage}" alt="Profile" class="w-8 h-8 rounded-full object-cover shadow-lg ring-2 ring-slate-400/40">
                        ` : `
                            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 via-slate-400 to-slate-500 flex items-center justify-center text-white shadow-lg ring-2 ring-slate-400/40">
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                                </svg>
                            </div>
                        `}
                    </div>
                    <div class="flex-1 text-center md:text-left">
                        <h2 class="text-2xl font-extrabold gradient-text mb-2 tracking-tight">${escapeHtml(profile.username)}</h2>
                        <div class="flex items-center justify-center md:justify-start gap-2 mb-2 text-gray-300 text-sm">
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
                            </svg>
                            <span>${escapeHtml(profile.email)}</span>
                        </div>
                        ${profile.bio ? `
                            <p class="text-gray-300 mb-3 leading-relaxed">${escapeHtml(profile.bio)}</p>
                        ` : `
                            <p class="text-gray-400 italic mb-3">No bio yet</p>
                        `}
                        <div class="flex items-center justify-center md:justify-start gap-2 text-sm text-gray-400">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                            </svg>
                            <span>Member since ${joinDate}</span>
                        </div>
                    </div>
                </div>
            </div>
            
            <!-- Stats -->
            <div class="grid grid-cols-2 md:grid-cols-4 gap-3 mb-6">
                <div class="card text-center p-4">
                    <div class="text-2xl font-bold text-gray-50 mb-1">${profile.postCount || 0}</div>
                    <div class="text-xs font-semibold text-gray-400">Posts</div>
                </div>
                <div class="card text-center p-4">
                    <div class="text-2xl font-bold text-gray-50 mb-1">${profile.commentCount || 0}</div>
                    <div class="text-xs font-semibold text-gray-400">Comments</div>
                </div>
                <div class="card text-center p-4">
                    <div class="text-2xl font-bold text-gray-50 mb-1">${profile.likesReceived || 0}</div>
                    <div class="text-xs font-semibold text-gray-400">Likes Received</div>
                </div>
                <div class="card text-center p-4">
                    <div class="text-2xl font-bold text-gray-50 mb-1">${profile.dislikesReceived || 0}</div>
                    <div class="text-xs font-semibold text-gray-400">Dislikes Received</div>
                </div>
            </div>
            
            <!-- Tabs -->
            <div class="mb-6">
                <div class="flex space-x-2 overflow-x-auto pb-2" style="scrollbar-width: none;">
                    <button onclick="loadProfileTab('overview')" id="tab-overview" class="category-filter-btn active whitespace-nowrap">
                        Overview
                    </button>
                    <button onclick="loadProfileTab('posts')" id="tab-posts" class="category-filter-btn whitespace-nowrap">
                        Posts (${profile.postCount || 0})
                    </button>
                    <button onclick="loadProfileTab('comments')" id="tab-comments" class="category-filter-btn whitespace-nowrap">
                        Comments (${profile.commentCount || 0})
                    </button>
                    <button onclick="loadProfileTab('likes')" id="tab-likes" class="category-filter-btn whitespace-nowrap">
                        Liked Posts
                    </button>
                </div>
            </div>
            
            <!-- Tab Content -->
            <div id="profile-tab-content">
                <div id="content-overview" class="profile-tab-panel">
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div class="card p-5">
                            <h3 class="text-lg font-bold text-gray-50 mb-4">Activity Summary</h3>
                            <div class="space-y-3">
                                <div class="flex justify-between items-center">
                                    <span class="text-gray-300">Likes Given</span>
                                    <span class="text-gray-50 font-bold">${profile.likesGiven || 0}</span>
                                </div>
                                <div class="flex justify-between items-center">
                                    <span class="text-gray-300">Dislikes Given</span>
                                    <span class="text-gray-50 font-bold">${profile.dislikesGiven || 0}</span>
                                </div>
                            </div>
                        </div>
                        <div class="card p-5">
                            <h3 class="text-lg font-bold text-gray-50 mb-4">Account Information</h3>
                            <div class="space-y-3">
                                <div class="flex justify-between items-center">
                                    <span class="text-gray-300">User ID</span>
                                    <span class="text-gray-50 font-bold">#${profile.userId}</span>
                                </div>
                                <div class="flex justify-between items-center">
                                    <span class="text-gray-300">Join Date</span>
                                    <span class="text-gray-50 font-bold">${joinDate}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                
                <div id="content-posts" class="profile-tab-panel hidden">
                    <div class="text-center py-8">
                        <div class="animate-spin rounded-full h-8 w-8 border-4 border-slate-300 border-t-slate-500 mx-auto mb-3"></div>
                        <p class="text-gray-500 text-sm">Loading posts...</p>
                    </div>
                </div>
                
                <div id="content-comments" class="profile-tab-panel hidden">
                    <div class="text-center py-8">
                        <div class="animate-spin rounded-full h-8 w-8 border-4 border-slate-300 border-t-slate-500 mx-auto mb-3"></div>
                        <p class="text-gray-500 text-sm">Loading comments...</p>
                    </div>
                </div>
                
                <div id="content-likes" class="profile-tab-panel hidden">
                    <div class="text-center py-8">
                        <div class="animate-spin rounded-full h-8 w-8 border-4 border-slate-300 border-t-slate-500 mx-auto mb-3"></div>
                        <p class="text-gray-500 text-sm">Loading liked posts...</p>
                    </div>
                </div>
            </div>
        </div>
    `;

    // Store profile data for tab loading
    window.currentProfile = profile;

    // Ensure overview tab is active and visible on initial load
    setTimeout(() => {
        loadProfileTab('overview');
    }, 0);
}

// Profile tab loading function
window.loadProfileTab = async function loadProfileTab(tabName) {
    // Update tab buttons
    document.querySelectorAll('.category-filter-btn').forEach(btn => {
        if (btn.id && btn.id.startsWith('tab-')) {
            btn.classList.remove('active');
        }
    });

    const activeTab = document.getElementById(`tab-${tabName}`);
    if (activeTab) {
        activeTab.classList.add('active');
    }

    // Hide all tab panels
    document.querySelectorAll('.profile-tab-panel').forEach(panel => {
        panel.classList.add('hidden');
    });

    // Show selected tab panel
    const contentPanel = document.getElementById(`content-${tabName}`);
    if (contentPanel) {
        contentPanel.classList.remove('hidden');
    }

    // Load content based on tab
    switch (tabName) {
        case 'posts':
            await loadUserPosts();
            break;
        case 'comments':
            await loadUserComments();
            break;
        case 'likes':
            await loadUserLikes();
            break;
        case 'overview':
            // Overview is already loaded
            break;
    }
}

// Load user posts
async function loadUserPosts() {
    const contentPanel = document.getElementById('content-posts');
    if (!contentPanel) return;

    try {
        const response = await fetch('/api/user/posts', {
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Failed to load posts');
        }

        const posts = await response.json();

        if (posts.length === 0) {
            contentPanel.innerHTML = `
                <div class="card text-center py-8">
                    <div class="text-4xl mb-3 animate-bounce">📝</div>
                    <h3 class="text-lg font-bold text-gray-100 mb-2">No posts yet</h3>
                    <p class="text-gray-400 mb-4 text-sm">Start sharing your thoughts with the community!</p>
                    <a href="#" onclick="navigateTo('/create-post'); return false;" class="btn-primary inline-flex items-center gap-2">
                        <span>✨</span> Create First Post
                    </a>
                </div>
            `;
            return;
        }

        // Use the same renderPosts function for consistency
        const postsContainer = document.createElement('div');
        postsContainer.id = 'posts-container';
        postsContainer.className = 'space-y-4';
        contentPanel.innerHTML = '';
        contentPanel.appendChild(postsContainer);
        renderPosts(posts);
    } catch (error) {
        console.error('Failed to load posts:', error);
        contentPanel.innerHTML = `
            <div class="card text-center py-8">
                <div class="text-2xl mb-3">😕</div>
                <p class="text-red-500 font-semibold mb-2">Failed to load posts</p>
                <p class="text-sm text-gray-500 mb-4">${error.message || 'Please try again later'}</p>
                <button onclick="loadProfileTab('posts')" class="btn-primary">🔄 Retry</button>
            </div>
        `;
    }
}

// Load user comments
async function loadUserComments() {
    const contentPanel = document.getElementById('content-comments');
    if (!contentPanel) return;

    try {
        const response = await fetch('/api/user/comments', {
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Failed to load comments');
        }

        const comments = await response.json();

        if (comments.length === 0) {
            contentPanel.innerHTML = `
                <div class="card text-center py-8">
                    <div class="text-4xl mb-3 animate-bounce">💬</div>
                    <h3 class="text-lg font-bold text-gray-100 mb-2">No comments yet</h3>
                    <p class="text-gray-400 mb-4 text-sm">Engage with the community by commenting on posts!</p>
                </div>
            `;
            return;
        }

        contentPanel.innerHTML = `
            <div class="space-y-4">
                ${comments.map((comment, index) => `
                    <article class="post-card group" style="animation: fadeInUp 0.5s ease-out ${index * 0.06}s both;">
                        <div class="post-card-content">
                            <div class="mb-3">
                                <a href="#" onclick="navigateTo('/post/${comment.postId}'); return false;" class="text-lg font-bold text-gray-50 mb-2 group-hover:text-slate-300 transition-colors leading-tight">
                                    ${escapeHtml(comment.title || 'Untitled Post')}
                                </a>
                            </div>
                            <p class="text-gray-300 mb-3 leading-relaxed text-sm">${escapeHtml(comment.content)}</p>
                            <div class="flex items-center gap-1.5 text-xs text-gray-400">
                                <svg class="w-3 h-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                                </svg>
                                <span>${comment.timeAgo || 'Recently'}</span>
                            </div>
                        </div>
                    </article>
                `).join('')}
            </div>
        `;
    } catch (error) {
        console.error('Failed to load comments:', error);
        contentPanel.innerHTML = `
            <div class="card text-center py-8">
                <div class="text-2xl mb-3">😕</div>
                <p class="text-red-500 font-semibold mb-2">Failed to load comments</p>
                <p class="text-sm text-gray-500 mb-4">${error.message || 'Please try again later'}</p>
                <button onclick="loadProfileTab('comments')" class="btn-primary">🔄 Retry</button>
            </div>
        `;
    }
}

// Load user liked posts
async function loadUserLikes() {
    const contentPanel = document.getElementById('content-likes');
    if (!contentPanel) return;

    try {
        const response = await fetch('/api/user/likes', {
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Failed to load liked posts');
        }

        const posts = await response.json();

        if (posts.length === 0) {
            contentPanel.innerHTML = `
                <div class="card text-center py-16">
                    <div class="text-3xl mb-6 animate-bounce">👍</div>
                    <h3 class="text-2xl font-bold text-gray-100 mb-2">No liked posts yet</h3>
                    <p class="text-gray-400 mb-6">Like posts you find interesting!</p>
                </div>
            `;
            return;
        }

        // Use the same renderPosts function for consistency
        const postsContainer = document.createElement('div');
        postsContainer.id = 'posts-container';
        postsContainer.className = 'space-y-4';
        contentPanel.innerHTML = '';
        contentPanel.appendChild(postsContainer);
        renderPosts(posts);
    } catch (error) {
        console.error('Failed to load liked posts:', error);
        contentPanel.innerHTML = `
            <div class="card text-center py-8">
                <div class="text-4xl mb-3">😕</div>
                <p class="text-red-500 font-semibold mb-2">Failed to load liked posts</p>
                <p class="text-sm text-gray-500 mb-4">${error.message || 'Please try again later'}</p>
                <button onclick="loadProfileTab('likes')" class="btn-primary">🔄 Retry</button>
            </div>
        `;
    }
}

// Render create post page
async function renderCreatePost() {
    const app = document.getElementById('app');

    // Load categories
    let categories = [];
    try {
        const response = await fetch('/api/categories');
        categories = await response.json();
    } catch (error) {
        console.error('Failed to load categories:', error);
    }

    app.innerHTML = `
        <div class="max-w-3xl mx-auto">
            <h1 class="text-3xl font-bold gradient-text mb-6">Create New Post</h1>
            <div class="card">
                <form id="create-post-form" class="space-y-6">
                    <div>
                        <label class="block text-gray-700 font-semibold mb-2">Title</label>
                        <input type="text" name="title" required class="input-field" placeholder="Enter post title">
                    </div>
                    <div>
                        <label class="block text-gray-700 font-semibold mb-2">Content</label>
                        <textarea name="content" required rows="10" class="input-field" placeholder="Write your post content..."></textarea>
                    </div>
                    <div>
                        <label class="block text-gray-700 font-semibold mb-2">Categories</label>
                        <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
                            ${categories.map(cat => `
                                <label class="flex items-center space-x-2 cursor-pointer">
                                    <input type="checkbox" name="categories[]" value="${escapeHtml(cat.name)}" class="w-5 h-5 text-slate-400 rounded">
                                    <span class="text-gray-700">${escapeHtml(cat.name)}</span>
                                </label>
                            `).join('')}
                        </div>
                    </div>
                    <div id="create-post-error" class="text-red-500 text-sm hidden"></div>
                    <button type="submit" class="btn-primary w-full">Create Post</button>
                </form>
            </div>
        </div>
    `;

    document.getElementById('create-post-form').addEventListener('submit', handleCreatePost);
}

// Handle create post
async function handleCreatePost(e) {
    e.preventDefault();
    const formData = new FormData(e.target);
    const errorDiv = document.getElementById('create-post-error');

    try {
        const response = await fetch('/new-post', {
            method: 'POST',
            body: formData
        });

        if (response.ok || response.redirected) {
            const location = response.headers.get('Location');
            if (location) {
                const postId = location.split('id=')[1];
                navigateTo(`/post/${postId}`);
            } else {
                navigateTo('/');
            }
        } else {
            errorDiv.textContent = 'Failed to create post. Please try again.';
            errorDiv.classList.remove('hidden');
        }
    } catch (error) {
        errorDiv.textContent = 'Failed to create post. Please try again.';
        errorDiv.classList.remove('hidden');
    }
}

// Users status box functions
function showUsersStatusBox() {
    const statusBox = document.getElementById('users-status-box');
    const app = document.getElementById('app');
    const chatInterface = document.getElementById('chat-interface');
    const usersList = document.getElementById('users-status-list');

    if (statusBox) {
        statusBox.classList.remove('hidden');
        statusBox.style.display = 'flex'; // Force display
        statusBox.style.visibility = 'visible'; // Force visibility
        statusBox.style.opacity = '1'; // Force opacity

        // Ensure chat interface is hidden
        if (chatInterface) {
            chatInterface.classList.add('hidden');
            chatInterface.style.display = 'none';
        }

        // Ensure users list is visible
        if (usersList) {
            usersList.style.display = 'block';
            usersList.style.visibility = 'visible';
        }

        if (app) {
            app.classList.add('with-sidebar');
        }
        console.log('Users status box shown', statusBox);
        console.log('Box position:', statusBox.getBoundingClientRect());
    } else {
        console.error('Users status box element not found in DOM!');
    }
}

function hideUsersStatusBox() {
    const statusBox = document.getElementById('users-status-box');
    const app = document.getElementById('app');
    if (statusBox) {
        statusBox.classList.add('hidden');
        statusBox.style.display = 'none';
        statusBox.style.visibility = 'hidden';
        if (app) {
            app.classList.remove('with-sidebar');
        }
    }
}

// Load users status
async function loadUsersStatus() {
    if (!AppState.currentUser) {
        console.log('No current user, skipping loadUsersStatus');
        return;
    }

    try {
        console.log('Loading users status...');
        const response = await fetch('/api/messages/users', {
            credentials: 'include'
        });

        console.log('Response status:', response.status);

        if (!response.ok) {
            const errorText = await response.text();
            console.error('HTTP error response:', errorText);
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const users = await response.json();
        console.log('Received users:', users);
        console.log('Users count:', users ? users.length : 0);

        // Handle null or undefined response
        if (!users || !Array.isArray(users)) {
            console.warn('Users list is not an array:', users);
            renderUsersStatus([]);
            return;
        }

        if (users.length === 0) {
            console.log('No users found (empty array returned) - this is normal if you are the only user');
        }

        // Initialize unread message counts from backend response
        users.forEach(user => {
            if (user.id && user.unread_count !== undefined) {
                AppState.unreadMessages[user.id] = user.unread_count || 0;
            }
        });

        renderUsersStatus(users);
    } catch (error) {
        console.error('Failed to load users status:', error);
        // Render empty list on error
        renderUsersStatus([]);
    }
}

// Render users status
function renderUsersStatus(users) {
    const container = document.getElementById('users-status-list');
    if (!container) {
        console.error('users-status-list container not found!');
        return;
    }

    // Make sure the parent box is visible
    const statusBox = document.getElementById('users-status-box');
    if (statusBox && statusBox.classList.contains('hidden')) {
        console.log('Status box was hidden, showing it now...');
        statusBox.classList.remove('hidden');
    }

    // Handle null, undefined, or non-array
    if (!users || !Array.isArray(users)) {
        console.warn('renderUsersStatus received invalid data:', users);
        container.innerHTML = '<div class="text-center py-12"><p class="text-gray-400 text-sm">No users found</p></div>';
        return;
    }

    if (users.length === 0) {
        container.innerHTML = '<div class="text-center py-12"><p class="text-gray-400 text-sm">No users found</p></div>';
        return;
    }

    // Separate online and offline users
    const onlineUsers = users.filter(user => {
        const isOnline = user.is_online === true || user.is_online === 1 || user.is_online === 'true' || user.is_online === '1';
        return isOnline;
    });
    const offlineUsers = users.filter(user => {
        const isOnline = user.is_online === true || user.is_online === 1 || user.is_online === 'true' || user.is_online === '1';
        return !isOnline;
    });

    // Sort: online users first, then offline users
    const sortedUsers = [...onlineUsers, ...offlineUsers];

    let userIndex = 0;
    container.innerHTML = `
        ${onlineUsers.length > 0 ? `
            <div style="margin-bottom: 16px;">
                <div style="display: flex; align-items: center; gap: 8px; padding: 8px 12px; background: rgba(34, 197, 94, 0.1); border-radius: 12px; border: 1px solid rgba(34, 197, 94, 0.2);">
                    <div style="width: 8px; height: 8px; border-radius: 50%; background: #22c55e; box-shadow: 0 0 12px rgba(34, 197, 94, 0.6); animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;"></div>
                    <span style="font-size: 11px; font-weight: 700; color: #4ade80; text-transform: uppercase; letter-spacing: 0.5px;">Online • ${onlineUsers.length}</span>
                </div>
            </div>
            <div style="display: flex; flex-direction: column; gap: 8px; margin-bottom: 20px;">
                ${onlineUsers.map((user) => {
        const index = userIndex++;
        const isOnline = true; // This user is online
        const unreadCount = AppState.unreadMessages[user.id] || 0;
        return `
                    <div class="user-status-item" data-user-id="${user.id}" onclick="openChat(${user.id}, '${escapeHtml(user.username)}', ${isOnline})" style="display: flex; align-items: center; gap: 12px; padding: 12px; border-radius: 14px; transition: all 0.2s ease; background: rgba(30, 41, 59, 0.4); border: 1px solid rgba(148, 163, 184, 0.1); cursor: pointer; animation: fadeInSlide 0.3s ease-out ${index * 0.03}s both; position: relative;" 
                         onmouseover="this.style.background='rgba(30, 41, 59, 0.7)'; this.style.borderColor='rgba(34, 197, 94, 0.3)'; this.style.transform='translateX(4px)'"
                         onmouseout="this.style.background='rgba(30, 41, 59, 0.4)'; this.style.borderColor='rgba(148, 163, 184, 0.1)'; this.style.transform='translateX(0)'">
                        <div style="position: relative; flex-shrink: 0;">
                            <div style="width: 36px; height: 36px; border-radius: 50%; background: linear-gradient(135deg, #22c55e 0%, #10b981 50%, #059669 100%); display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 12px rgba(34, 197, 94, 0.3), 0 0 0 2px rgba(34, 197, 94, 0.2);">
                                <svg style="width: 18px; height: 18px; color: white;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                    </svg>
                </div>
                            <div style="position: absolute; bottom: -2px; right: -2px; width: 14px; height: 14px; background: #22c55e; border-radius: 50%; border: 3px solid #0f172a; box-shadow: 0 0 0 2px rgba(34, 197, 94, 0.5);">
                                <div style="position: absolute; inset: 0; background: #22c55e; border-radius: 50%; animation: ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite; opacity: 0.75;"></div>
                        </div>
                    </div>
                        <div style="flex: 1; min-width: 0;">
                            <div style="display: flex; align-items: center; gap: 10px;">
                                <div style="position: relative; flex-shrink: 0;" title="Online">
                                    <div style="width: 10px; height: 10px; border-radius: 50%; background: #22c55e; box-shadow: 0 0 8px rgba(34, 197, 94, 0.6), 0 0 0 2px rgba(34, 197, 94, 0.2);"></div>
                                    <div style="position: absolute; inset: 0; width: 10px; height: 10px; border-radius: 50%; background: #22c55e; animation: ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite; opacity: 0.6;"></div>
            </div>
                                <h4 style="font-weight: 600; color: #f1f5f9; font-size: 14px; margin: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">${escapeHtml(user.username)}</h4>
                            </div>
                        </div>
                        ${unreadCount > 0 ? `
                            <div style="flex-shrink: 0; position: relative;">
                                <div style="min-width: 20px; height: 20px; padding: 0 6px; background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); color: white; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 700; box-shadow: 0 2px 8px rgba(239, 68, 68, 0.4);">
                                    ${unreadCount > 99 ? '99+' : unreadCount}
                                </div>
                            </div>
                        ` : ''}
                    </div>
        `;
    }).join('')}
                </div>
        ` : ''}
        ${offlineUsers.length > 0 ? `
            <div style="margin-top: ${onlineUsers.length > 0 ? '24px' : '0'}; margin-bottom: 16px; padding-top: ${onlineUsers.length > 0 ? '20px' : '0'}; border-top: ${onlineUsers.length > 0 ? '1px solid rgba(148, 163, 184, 0.1)' : 'none'};">
                <div style="display: flex; align-items: center; gap: 8px; padding: 8px 12px; background: rgba(100, 116, 139, 0.1); border-radius: 12px; border: 1px solid rgba(100, 116, 139, 0.15);">
                    <div style="width: 8px; height: 8px; border-radius: 50%; background: #64748b;"></div>
                    <span style="font-size: 11px; font-weight: 700; color: #94a3b8; text-transform: uppercase; letter-spacing: 0.5px;">Offline • ${offlineUsers.length}</span>
                </div>
                <div style="margin-top: 8px; padding: 6px 12px;">
                    <p style="font-size: 10px; color: #64748b; margin: 0; font-style: italic;">Recently online (last hour)</p>
                </div>
            </div>
        ` : ''}
        ${offlineUsers.length > 0 ? `
            <div style="display: flex; flex-direction: column; gap: 8px;">
                ${offlineUsers.map((user) => {
        const index = userIndex++;
        const isOffline = false; // This user is offline
        const unreadCount = AppState.unreadMessages[user.id] || 0;
        return `
                    <div class="user-status-item" data-user-id="${user.id}" onclick="openChat(${user.id}, '${escapeHtml(user.username)}', ${isOffline})" style="display: flex; align-items: center; gap: 12px; padding: 12px; border-radius: 14px; transition: all 0.2s ease; background: rgba(30, 41, 59, 0.3); border: 1px solid rgba(148, 163, 184, 0.08); cursor: pointer; animation: fadeInSlide 0.3s ease-out ${index * 0.03}s both; opacity: 0.7; position: relative;" 
                         onmouseover="this.style.background='rgba(30, 41, 59, 0.5)'; this.style.borderColor='rgba(148, 163, 184, 0.2)'; this.style.opacity='1'; this.style.transform='translateX(4px)'"
                         onmouseout="this.style.background='rgba(30, 41, 59, 0.3)'; this.style.borderColor='rgba(148, 163, 184, 0.08)'; this.style.opacity='0.7'; this.style.transform='translateX(0)'">
                        <div style="position: relative; flex-shrink: 0;">
                            <div style="width: 44px; height: 44px; border-radius: 50%; background: linear-gradient(135deg, #64748b 0%, #475569 50%, #334155 100%); display: flex; align-items: center; justify-content: center; box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);">
                                <svg style="width: 22px; height: 22px; color: white;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
        </svg>
                            </div>
                            <div style="position: absolute; bottom: -2px; right: -2px; width: 14px; height: 14px; background: #64748b; border-radius: 50%; border: 3px solid #0f172a;"></div>
                        </div>
                        <div style="flex: 1; min-width: 0;">
                            <div style="display: flex; align-items: center; gap: 10px;">
                                <div style="width: 10px; height: 10px; border-radius: 50%; background: #64748b; flex-shrink: 0;" title="Offline"></div>
                                <h4 style="font-weight: 600; color: #94a3b8; font-size: 14px; margin: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">${escapeHtml(user.username)}</h4>
                            </div>
                        </div>
                        ${unreadCount > 0 ? `
                            <div style="flex-shrink: 0; position: relative;">
                                <div style="min-width: 20px; height: 20px; padding: 0 6px; background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); color: white; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 700; box-shadow: 0 2px 8px rgba(239, 68, 68, 0.4);">
                                    ${unreadCount > 99 ? '99+' : unreadCount}
                                </div>
                            </div>
                        ` : ''}
                    </div>
                `;
    }).join('')}
                </div>
        ` : ''}
        ${users.length === 0 ? `
            <div style="text-align: center; padding: 48px 24px;">
                <div style="width: 64px; height: 64px; margin: 0 auto 16px; background: rgba(148, 163, 184, 0.1); border-radius: 50%; display: flex; align-items: center; justify-content: center;">
                    <svg style="width: 32px; height: 32px; color: #64748b;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
                    </svg>
                </div>
                <p style="color: #94a3b8; font-size: 14px; margin: 0 0 8px 0; font-weight: 500;">No other users found</p>
                <p style="color: #64748b; font-size: 12px; margin: 0;">Other users will appear here when they register</p>
            </div>
        ` : ''}
    `;
}


// Open chat with user
function openChat(userId, username, isOnline) {
    AppState.currentChatUser = { id: userId, username, isOnline };
    AppState.messagesOffset = 0;
    AppState.hasMoreMessages = true;

    // Clear unread count for this user
    AppState.unreadMessages[userId] = 0;

    // Hide users list and show chat interface
    const usersList = document.getElementById('users-status-list');
    const chatInterface = document.getElementById('chat-interface');

    if (usersList) {
        usersList.style.display = 'none';
    }

    if (chatInterface) {
        chatInterface.classList.remove('hidden');
        chatInterface.style.display = 'flex';
    }

    // Refresh users list to update unread badges
    loadUsersStatus();

    // Update chat header
    document.getElementById('chat-username').textContent = username;
    const statusEl = document.getElementById('chat-status');
    const avatarEl = document.getElementById('chat-user-avatar');

    if (isOnline) {
        statusEl.innerHTML = `
            <div style="display: flex; align-items: center; gap: 6px;">
                <div style="position: relative;">
                    <div style="width: 8px; height: 8px; background: #22c55e; border-radius: 50%; box-shadow: 0 0 8px rgba(34, 197, 94, 0.6);"></div>
                    <div style="position: absolute; inset: 0; width: 8px; height: 8px; background: #22c55e; border-radius: 50%; animation: ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite; opacity: 0.75;"></div>
                </div>
                <span style="font-size: 11px; font-weight: 600; color: #4ade80; text-transform: uppercase; letter-spacing: 0.5px;">Online</span>
            </div>
        `;
        avatarEl.style.background = 'linear-gradient(135deg, #22c55e 0%, #10b981 50%, #059669 100%)';
        avatarEl.style.boxShadow = '0 4px 12px rgba(34, 197, 94, 0.3), 0 0 0 2px rgba(34, 197, 94, 0.2)';
    } else {
        statusEl.innerHTML = `
            <div style="display: flex; align-items: center; gap: 6px;">
                <div style="width: 8px; height: 8px; background: #64748b; border-radius: 50%;"></div>
                <span style="font-size: 11px; font-weight: 600; color: #94a3b8; text-transform: uppercase; letter-spacing: 0.5px;">Offline</span>
            </div>
        `;
        avatarEl.style.background = 'linear-gradient(135deg, #64748b 0%, #475569 50%, #334155 100%)';
        avatarEl.style.boxShadow = '0 2px 8px rgba(0, 0, 0, 0.2)';
    }

    // Setup back button
    const backBtn = document.getElementById('back-to-users-btn');
    if (backBtn) {
        // Remove any existing event listeners by cloning
        const newBtn = backBtn.cloneNode(true);
        backBtn.parentNode.replaceChild(newBtn, backBtn);

        // Add new event listener
        const backButton = document.getElementById('back-to-users-btn');
        if (backButton) {
            backButton.onclick = (e) => {
                e.preventDefault();
                e.stopPropagation();
                console.log('Back button clicked, returning to users list');

                // Hide chat interface first
                const chatInterface = document.getElementById('chat-interface');
                if (chatInterface) {
                    chatInterface.classList.add('hidden');
                    chatInterface.style.display = 'none';
                }

                // Show users list
                const usersList = document.getElementById('users-status-list');
                if (usersList) {
                    usersList.style.display = 'block';
                    usersList.style.visibility = 'visible';
                }

                AppState.currentChatUser = null;
                loadUsersStatus();
            };
        }
    }

    // Setup form
    const chatForm = document.getElementById('chat-form');
    const chatInput = document.getElementById('chat-input');

    // Remove any existing event listeners by cloning the input
    const newInput = chatInput.cloneNode(true);
    chatInput.parentNode.replaceChild(newInput, chatInput);
    const chatInputNew = document.getElementById('chat-input');

    // Typing indicator logic
    let typingTimeout = null;
    let isCurrentlyTyping = false;

    // Typing detection - send immediately when user starts typing
    const handleTyping = () => {
        if (!isCurrentlyTyping) {
            isCurrentlyTyping = true;
            sendTypingIndicator(userId, true);
        }

        // Clear existing timeout
        if (typingTimeout) {
            clearTimeout(typingTimeout);
        }

        // Set timeout to send typing_stop after user stops typing
        typingTimeout = setTimeout(() => {
            if (isCurrentlyTyping) {
                isCurrentlyTyping = false;
                sendTypingIndicator(userId, false);
            }
        }, AppState.typingDebounceTime);
    };

    // Auto-resize textarea
    const autoResize = () => {
        chatInputNew.style.height = 'auto';
        chatInputNew.style.height = Math.min(chatInputNew.scrollHeight, 120) + 'px';
    };

    // Listen for input events (for auto-resize and typing detection)
    chatInputNew.addEventListener('input', (e) => {
        autoResize();
        handleTyping();
    });

    // Listen for keydown events (for better typing responsiveness and Enter handling)
    chatInputNew.addEventListener('keydown', (e) => {
        // Handle Enter key for sending
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            // Stop typing indicator
            if (isCurrentlyTyping) {
                isCurrentlyTyping = false;
                sendTypingIndicator(userId, false);
            }
            if (typingTimeout) {
                clearTimeout(typingTimeout);
            }
            sendMessage(userId);
            return;
        }
        
        // Trigger typing detection for character keys, backspace, delete
        if (e.key.length === 1 || e.key === 'Backspace' || e.key === 'Delete') {
            handleTyping();
        }
    });

    // Stop typing when input loses focus
    chatInputNew.addEventListener('blur', () => {
        if (isCurrentlyTyping) {
            isCurrentlyTyping = false;
            sendTypingIndicator(userId, false);
        }
        if (typingTimeout) {
            clearTimeout(typingTimeout);
        }
    });

    // Form submit handler
    chatForm.onsubmit = (e) => {
        e.preventDefault();
        // Stop typing indicator
        if (isCurrentlyTyping) {
            isCurrentlyTyping = false;
            sendTypingIndicator(userId, false);
        }
        if (typingTimeout) {
            clearTimeout(typingTimeout);
        }
        sendMessage(userId);
    };

    // Focus input when chat opens
    setTimeout(() => {
        chatInputNew.focus();
    }, 100);

    // Setup "Load More" button handler
    setupLoadMoreButton(userId);

    // Load messages
    loadMessages(userId);
}

// Load messages
async function loadMessages(userId, append = false) {
    if (AppState.isLoadingMessages) return;
    AppState.isLoadingMessages = true;

    // Update button to show loading state
    if (AppState.updateLoadMoreButton) {
        AppState.updateLoadMoreButton();
    }

    const container = document.getElementById('chat-messages');
    let previousScrollHeight = 0;
    let previousScrollTop = 0;

    // Save scroll position before loading (for append mode)
    if (append && container) {
        previousScrollHeight = container.scrollHeight;
        previousScrollTop = container.scrollTop;
    }

    try {
        const response = await fetch(`/api/messages?user_id=${userId}&offset=${AppState.messagesOffset}&limit=10`, {
            credentials: 'include'
        });

        // Handle different response statuses
        if (response.status === 401 || response.status === 403) {
            // Authentication errors - show error
            throw new Error('Unauthorized');
        } else if (response.status >= 500) {
            // Server errors - show error
            throw new Error('Server error');
        } else if (!response.ok && response.status >= 400) {
            // Other client errors - but for 404 or empty responses, treat as no messages
            if (response.status === 404) {
                // 404 means no messages exist - not an error
                renderMessages([], append);
                AppState.hasMoreMessages = false;
                AppState.isLoadingMessages = false;
                return;
            }
            throw new Error('Failed to fetch messages');
        }

        const messages = await response.json();

        if (!Array.isArray(messages)) {
            throw new Error('Invalid messages format');
        }

        // If messages array is empty or has fewer than 10 messages, no more messages available
        if (messages.length < 10) {
            AppState.hasMoreMessages = false;
        }

        renderMessages(messages, append, previousScrollHeight, previousScrollTop);
        AppState.messagesOffset += messages.length;

        // Scroll to bottom if not appending (new chat)
        if (!append) {
            if (container) {
                setTimeout(() => {
                    container.scrollTop = container.scrollHeight;
                }, 100);
            }
        }
    } catch (error) {
        console.error('Failed to load messages:', error);
        if (container && !append) {
            // Check if it's a network error (no response) vs a server error
            // For network errors or actual server errors, show error message
            // For other cases (like 404 or empty), show empty state
            if (error.message === 'Unauthorized' || error.message === 'Server error') {
                container.innerHTML = `<div style="text-align: center; padding: 32px; color: #ef4444;">Failed to load messages. Please try again.</div>`;
            } else {
                // For any other error (including network errors, parsing errors, etc.)
                // If we can't load messages, assume it's because there are no messages
                // and show the empty state instead of an error
                renderMessages([], append);
            }
        }
    } finally {
        AppState.isLoadingMessages = false;
        // Update button state after loading completes
        if (AppState.updateLoadMoreButton) {
            AppState.updateLoadMoreButton();
        }
    }
}

// Setup "Load More" button
function setupLoadMoreButton(userId) {
    const container = document.getElementById('chat-messages');
    if (!container) return;

    // Remove existing button if any
    const existingButton = document.getElementById('load-more-messages-btn');
    if (existingButton) {
        existingButton.remove();
    }

    // Create throttled load more handler
    const throttledLoadMore = throttle(() => {
        if (!AppState.hasMoreMessages || AppState.isLoadingMessages) {
            return;
        }
        loadMessages(userId, true);
    }, 500); // Throttle to 500ms to prevent spam

    // Function to update button visibility
    const updateLoadMoreButton = () => {
        const existingBtn = document.getElementById('load-more-messages-btn');
        if (AppState.hasMoreMessages) {
            if (!existingBtn) {
                // Create and insert button at the top
                const loadMoreBtn = document.createElement('div');
                loadMoreBtn.id = 'load-more-messages-btn';
                loadMoreBtn.innerHTML = `
                    <div style="display: flex; justify-content: center; padding: 12px 0;">
                        <button id="load-more-btn" 
                                style="padding: 10px 20px; background: linear-gradient(135deg, rgba(30, 41, 59, 0.8) 0%, rgba(15, 23, 42, 0.8) 100%); 
                                       border: 1px solid rgba(148, 163, 184, 0.3); 
                                       border-radius: 12px; 
                                       color: #f1f5f9; 
                                       font-size: 13px; 
                                       font-weight: 600; 
                                       cursor: pointer; 
                                       transition: all 0.2s;
                                       display: flex;
                                       align-items: center;
                                       gap: 8px;"
                                onmouseover="if (!this.disabled) { this.style.background='rgba(30, 41, 59, 1)'; this.style.borderColor='rgba(148, 163, 184, 0.5)'; this.style.transform='scale(1.05)'; }"
                                onmouseout="if (!this.disabled) { this.style.background='linear-gradient(135deg, rgba(30, 41, 59, 0.8) 0%, rgba(15, 23, 42, 0.8) 100%)'; this.style.borderColor='rgba(148, 163, 184, 0.3)'; this.style.transform='scale(1)'; }">
                            <svg style="width: 16px; height: 16px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18"></path>
                            </svg>
                            <span id="load-more-text">Load More Messages</span>
                        </button>
                    </div>
                `;
                container.insertBefore(loadMoreBtn, container.firstChild);
                
                // Attach click handler
                const btn = document.getElementById('load-more-btn');
                if (btn) {
                    btn.addEventListener('click', throttledLoadMore);
                }
            } else {
                // Update button state
                const btn = document.getElementById('load-more-btn');
                const text = document.getElementById('load-more-text');
                if (btn && text) {
                    if (AppState.isLoadingMessages) {
                        btn.disabled = true;
                        btn.style.opacity = '0.6';
                        btn.style.cursor = 'not-allowed';
                        text.textContent = 'Loading...';
                    } else {
                        btn.disabled = false;
                        btn.style.opacity = '1';
                        btn.style.cursor = 'pointer';
                        text.textContent = 'Load More Messages';
                    }
                }
            }
        } else if (existingBtn) {
            existingBtn.remove();
        }
    };

    // Store update function for later use
    AppState.updateLoadMoreButton = updateLoadMoreButton;

    // Initial update
    updateLoadMoreButton();
}

// Render messages
function renderMessages(messages, append = false, previousScrollHeight = 0, previousScrollTop = 0) {
    const container = document.getElementById('chat-messages');
    if (!container) {
        console.error('chat-messages container not found');
        return;
    }

    // Save typing indicator state before clearing
    const typingIndicator = document.getElementById('typing-indicator');
    const wasTypingVisible = typingIndicator && typingIndicator.style.display !== 'none';
    const typingUsername = typingIndicator ? typingIndicator.querySelector('.typing-username')?.textContent : '';
    const typingIndicatorHTML = typingIndicator ? typingIndicator.outerHTML : '';

    // Save "Load More" button if it exists
    const loadMoreButton = document.getElementById('load-more-messages-btn');
    const loadMoreButtonHTML = loadMoreButton ? loadMoreButton.outerHTML : '';

    if (!append) {
        container.innerHTML = '';
    }

    const currentUserId = AppState.currentUser ? AppState.currentUser.id : null;
    if (!currentUserId) {
        console.error('No current user ID');
        return;
    }

    let lastDate = null;

    // If no messages, show empty state
    if (!messages || messages.length === 0) {
        if (!append) {
            container.innerHTML = `
                <div style="display: flex; align-items: center; justify-content: center; height: 100%; flex-direction: column; gap: 12px; opacity: 0.6;">
                    <svg style="width: 48px; height: 48px; color: #64748b;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
                    </svg>
                    <p style="color: #94a3b8; font-size: 14px; margin: 0;">Write something...</p>
                </div>
            `;
        }
        return;
    }

    const messagesHtml = messages.map((msg, index) => {
        const isSender = Number(msg.sender_id) === Number(currentUserId);

        // Format date
        const msgDate = msg.date || (msg.datetime ? msg.datetime.split(' ')[0] : null);
        const msgTime = msg.time || (msg.datetime ? msg.datetime.split(' ')[1] : null) || msg.time_ago || '';

        // Show date separator if date changed
        let dateSeparator = '';
        if (msgDate && msgDate !== lastDate) {
            lastDate = msgDate;
            const displayDate = new Date(msgDate).toLocaleDateString('en-US', {
                weekday: 'short',
                year: 'numeric',
                month: 'short',
                day: 'numeric'
            });
            dateSeparator = `
                <div style="display: flex; align-items: center; justify-content: center; margin: 20px 0 16px 0;">
                    <div style="padding: 6px 14px; background: rgba(30, 41, 59, 0.6); border-radius: 12px; border: 1px solid rgba(148, 163, 184, 0.2);">
                        <span style="font-size: 11px; color: #94a3b8; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px;">${displayDate}</span>
                    </div>
                </div>
            `;
        }

        return dateSeparator + `
            <div style="display: flex; ${isSender ? 'justify-content: flex-end' : 'justify-content: flex-start'}; margin-bottom: 6px; width: 100%; align-items: flex-end;">
                <div style="max-width: 75%; min-width: 60px; padding: 10px 14px; border-radius: 18px; ${isSender ? 'border-bottom-right-radius: 4px' : 'border-bottom-left-radius: 4px'}; ${isSender ? 'background: linear-gradient(135deg, #22c55e 0%, #10b981 100%); color: white; box-shadow: 0 2px 8px rgba(34, 197, 94, 0.3);' : 'background: rgba(30, 41, 59, 0.8); color: #f1f5f9; border: 1px solid rgba(148, 163, 184, 0.2); box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);'}; word-wrap: break-word; overflow-wrap: break-word; display: inline-block;">
                    <div style="font-size: 14px; line-height: 1.5; white-space: pre-wrap; word-break: break-word; margin-bottom: 4px; ${isSender ? 'color: white' : 'color: #f1f5f9'};">
                        ${escapeHtml(msg.content)}
                    </div>
                    <div style="font-size: 10px; ${isSender ? 'color: rgba(255, 255, 255, 0.8)' : 'color: #94a3b8'}; text-align: right; font-weight: 500; margin-top: 2px;">
                        ${msgTime}
                    </div>
                </div>
            </div>
        `;
    }).join('');

    if (append) {
        // Insert messages at the beginning (after the load more button if it exists)
        const loadMoreBtn = document.getElementById('load-more-messages-btn');
        if (loadMoreBtn) {
            loadMoreBtn.insertAdjacentHTML('afterend', messagesHtml);
        } else {
            container.insertAdjacentHTML('afterbegin', messagesHtml);
        }
        
        // Preserve scroll position after loading older messages
        if (previousScrollHeight > 0) {
            const newScrollHeight = container.scrollHeight;
            const scrollDifference = newScrollHeight - previousScrollHeight;
            container.scrollTop = previousScrollTop + scrollDifference;
        }
    } else {
        // Prepend load more button if there are more messages, then messages, then typing indicator
        let fullHTML = '';
        if (AppState.hasMoreMessages) {
            fullHTML = loadMoreButtonHTML;
        }
        fullHTML += messagesHtml + typingIndicatorHTML;
        container.innerHTML = fullHTML;

        // Restore typing indicator visibility if it was visible
        if (wasTypingVisible && typingIndicator) {
            const restoredTypingEl = document.getElementById('typing-indicator');
            if (restoredTypingEl) {
                restoredTypingEl.style.display = 'flex';
                const usernameEl = restoredTypingEl.querySelector('.typing-username');
                if (usernameEl && typingUsername) {
                    usernameEl.textContent = typingUsername;
                }
            }
        }
        
        // Re-setup load more button handler if needed
        if (AppState.currentChatUser) {
            const currentUserId = AppState.currentChatUser.id;
            // Use setTimeout to ensure DOM is updated first
            setTimeout(() => {
                setupLoadMoreButton(currentUserId);
            }, 10);
        }
        
        // Scroll to bottom after rendering
        setTimeout(() => {
            container.scrollTop = container.scrollHeight;
        }, 50);
    }
    
    // Update button state after rendering (for both append and non-append)
    if (AppState.updateLoadMoreButton) {
        AppState.updateLoadMoreButton();
    }
}

// Send message
async function sendMessage(receiverId) {
    const input = document.getElementById('chat-input');
    const submitBtn = document.querySelector('#chat-form button[type="submit"]');
    const content = input.value.trim();

    if (!content) return;

    // Disable input and button while sending
    input.disabled = true;
    submitBtn.disabled = true;

    const formData = new FormData();
    formData.append('receiver_id', receiverId);
    formData.append('content', content);

    try {
        const response = await fetch('/api/messages/send', {
            method: 'POST',
            body: formData,
            credentials: 'include'
        });

        if (response.ok) {
            input.value = '';
            // Reset textarea height
            input.style.height = 'auto';
            // Reload messages to show the new one
            AppState.messagesOffset = 0;
            await loadMessages(receiverId);
            await loadUsersStatus();
        } else {
            console.error('Failed to send message');
            alert('Failed to send message. Please try again.');
        }
    } catch (error) {
        console.error('Failed to send message:', error);
        alert('Failed to send message. Please try again.');
    } finally {
        // Re-enable input and button
        input.disabled = false;
        submitBtn.disabled = false;
        input.focus();
    }
}

// Make openChat available globally
window.openChat = openChat;

// Utility functions
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Proper throttle function (not debounce)
function throttle(func, wait) {
    let inThrottle;
    return function executedFunction(...args) {
        if (!inThrottle) {
            func.apply(this, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, wait);
        }
    };
}

// Setup event listeners
function setupEventListeners() {
    // Handle clicks on links with data-navigate attribute
    document.addEventListener('click', (e) => {
        const link = e.target.closest('a[data-navigate]');
        if (link) {
            e.preventDefault();
            const path = link.getAttribute('data-navigate');
            navigateTo(path);
        }
    });

    // Minimize status box button
    const minimizeBtn = document.getElementById('minimize-status');
    if (minimizeBtn) {
        minimizeBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            const statusBox = document.getElementById('users-status-box');
            const usersList = document.getElementById('users-status-list');
            const chatInterface = document.getElementById('chat-interface');

            if (statusBox) {
                const isMinimized = statusBox.classList.contains('minimized');

                if (isMinimized) {
                    // Expand: show full box
                    statusBox.classList.remove('minimized');
                    statusBox.style.height = '650px';
                    if (usersList) {
                        usersList.style.display = 'block';
                    }
                    if (chatInterface) {
                        chatInterface.style.display = chatInterface.classList.contains('hidden') ? 'none' : 'flex';
                    }
                    // Change button to X (minimize icon)
                    minimizeBtn.innerHTML = '×';
                    minimizeBtn.title = 'Minimize';
                } else {
                    // Minimize: show only header
                    statusBox.classList.add('minimized');
                    statusBox.style.height = 'auto';
                    if (usersList) {
                        usersList.style.display = 'none';
                    }
                    if (chatInterface) {
                        chatInterface.style.display = 'none';
                    }
                    // Change button to expand icon
                    minimizeBtn.innerHTML = '<svg style="width: 16px; height: 16px;" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4"></path></svg>';
                    minimizeBtn.title = 'Expand';
                }
            }
        });
    }
}

// Setup drag functionality for users status box
function setupUsersStatusBoxDrag() {
    const statusBox = document.getElementById('users-status-box');
    const header = document.getElementById('status-header');

    if (!statusBox || !header) return;

    let isDragging = false;
    let startX = 0;
    let startY = 0;
    let startRight = 0;
    let startTop = 0;

    header.addEventListener('mousedown', (e) => {
        if (e.target.id === 'minimize-status' || e.target.closest('#minimize-status')) return;

        isDragging = true;
        statusBox.classList.add('dragging');

        const rect = statusBox.getBoundingClientRect();
        startX = e.clientX;
        startY = e.clientY;
        startRight = window.innerWidth - rect.right;
        startTop = rect.top;

        e.preventDefault();
    });

    document.addEventListener('mousemove', (e) => {
        if (!isDragging) return;

        const deltaX = startX - e.clientX;
        const deltaY = e.clientY - startY;

        const newRight = Math.max(0, Math.min(window.innerWidth - 300, startRight + deltaX));
        const newTop = Math.max(64, Math.min(window.innerHeight - 500, startTop + deltaY));

        statusBox.style.right = newRight + 'px';
        statusBox.style.top = newTop + 'px';
    });

    document.addEventListener('mouseup', () => {
        if (isDragging) {
            isDragging = false;
            statusBox.classList.remove('dragging');
        }
    });
}

// Internal navigation function (doesn't update history)
function navigateInternal(path) {
    if (path === '/' || path === '') {
        AppState.currentPage = 'home';
        AppState.currentPostId = null;
    } else if (path === '/login') {
        AppState.currentPage = 'login';
        AppState.currentPostId = null;
    } else if (path === '/register') {
        AppState.currentPage = 'register';
        AppState.currentPostId = null;
    } else if (path.startsWith('/post/')) {
        const postId = path.split('/post/')[1];
        AppState.currentPage = 'post';
        AppState.currentPostId = postId;
    } else if (path === '/create-post') {
        AppState.currentPage = 'create-post';
        AppState.currentPostId = null;
    } else if (path === '/profile') {
        AppState.currentPage = 'profile';
        AppState.currentPostId = null;
    } else {
        AppState.currentPage = 'home';
        AppState.currentPostId = null;
    }

    renderPage();
}

// Make navigateTo available globally (updates history)
window.navigateTo = (path) => {
    window.history.pushState({}, '', path);
    navigateInternal(path);
};

// Add new functions for typing indicators
function handleTypingStart(data) {
    console.log('handleTypingStart called:', data);
    console.log('AppState.currentChatUser:', AppState.currentChatUser);

    // Only show if we're in a chat with this user
    if (!AppState.currentChatUser) return;

    const senderId = Number(data.sender_id);
    const currentChatUserId = Number(AppState.currentChatUser.id);

    console.log('senderId:', senderId, 'currentChatUserId:', currentChatUserId);

    if (senderId === currentChatUserId) {
        console.log('Showing typing indicator for:', data.sender_username);
        showTypingIndicator(data.sender_username || data.username || 'Someone');
    } else {
        console.log('Sender ID does not match current chat user ID');
    }
}

function handleTypingStop(data) {
    // Only hide if we're in a chat with this user
    if (!AppState.currentChatUser) return;

    const senderId = Number(data.sender_id);
    const currentChatUserId = Number(AppState.currentChatUser.id);

    if (senderId === currentChatUserId) {
        hideTypingIndicator();
    }
}

function showTypingIndicator(username) {
    console.log('showTypingIndicator called with username:', username);
    const typingEl = document.getElementById('typing-indicator');
    if (!typingEl) {
        console.error('Typing indicator element not found!');
        return;
    }

    // Clear any existing timeout
    if (AppState.typingTimeouts[username]) {
        clearTimeout(AppState.typingTimeouts[username]);
        delete AppState.typingTimeouts[username];
    }

    console.log('Setting typing indicator to visible');
    typingEl.style.display = 'flex';
    const usernameEl = typingEl.querySelector('.typing-username');
    if (usernameEl) {
        usernameEl.textContent = username;
    } else {
        console.error('Could not find .typing-username element');
    }

    // Scroll to bottom to show typing indicator
    const container = document.getElementById('chat-messages');
    if (container) {
        setTimeout(() => {
            container.scrollTop = container.scrollHeight;
        }, 50);
    }
}

function hideTypingIndicator() {
    const typingEl = document.getElementById('typing-indicator');
    if (typingEl) {
        typingEl.style.display = 'none';
    }
}

// Send typing indicator to server
function sendTypingIndicator(receiverId, isTyping) {
    if (!AppState.ws || AppState.ws.readyState !== WebSocket.OPEN) return;
    if (!AppState.currentUser) return;

    const message = {
        type: isTyping ? 'typing_start' : 'typing_stop',
        sender_id: AppState.currentUser.id,
        sender_username: AppState.currentUser.username,
        receiver_id: receiverId
    };

    AppState.ws.send(JSON.stringify(message));
}


