import { me } from './user_store.svelte';
import { api } from './api';

export const Permissions = [
	'administrator',
	'manage_server',
	'manage_roles',
	'manage_channels',
	'kick_users',
	'ban_users',
	'create_invites',
	'manage_invites',
	'send_messages',
	'manage_messages',
	'embed_links',
	'attach_files',
	'mention_everyone',
	'view_channels'
] as const;

export type Permission = (typeof Permissions)[number];

// cache for user permissions per server
const permissionCache: Record<string, string[]> = {};

function getPermissions(serverId: string): Permission[] {
	if (!me.user) return [];
	const server = me.servers?.find((s) => s.id === serverId);
	if (!server) return [];
	if (server.owner_id === me.user.id) {
		return [...Permissions];
	}
	return (permissionCache[serverId] as Permission[]) || ['send_messages', 'view_channels'];
}

// fetch and cache permissions
export async function fetchUserPermissions(serverId: string): Promise<string[]> {
	if (!me.user) return [];
	const server = me.servers?.find((s) => s.id === serverId);
	if (!server) return [];
	if (server.owner_id === me.user.id) {
		return getPermissions(serverId);
	}
	try {
		const roles = await api.server.userRoles(serverId);
		const perms = new Set(roles.flatMap((role) => role.permissions));
		const permArray = Array.from(perms);
		permissionCache[serverId] = permArray;
		return permArray;
	} catch (error) {
		console.error('failed to fetch user permissions:', error);
		return ['send_messages', 'view_channels'];
	}
}

// invalidate cache for a server
export function invalidatePermissionCache(serverId: string) {
	delete permissionCache[serverId];
}

// check if user has a specific permission in a server
export function can(permission: Permission, serverId: string): boolean {
	const perms = getPermissions(serverId);
	return perms.includes('administrator') || perms.includes(permission);
}
