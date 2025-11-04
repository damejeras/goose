import {useState, useEffect} from "react";
import {apiClient} from "@/lib/apiClient";
import type {APIKey} from "@/lib/api/v1/apikey_pb";
import {Heading} from "@/components/heading";
import {Button} from "@/components/button";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/table";
import {Badge} from "@/components/badge";
import {Dialog, DialogActions, DialogBody, DialogDescription, DialogTitle} from "@/components/dialog";
import {Field, Label} from "@/components/fieldset";
import {Input} from "@/components/input";
import {TrashIcon, PencilIcon, PlusIcon} from "@heroicons/react/20/solid";
import {Text} from "@/components/text";

export default function APIKeys() {
    const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Create dialog state
    const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
    const [createName, setCreateName] = useState("");
    const [createdKey, setCreatedKey] = useState<string | null>(null);
    const [isCreating, setIsCreating] = useState(false);

    // Edit dialog state
    const [editingKey, setEditingKey] = useState<APIKey | null>(null);
    const [editName, setEditName] = useState("");
    const [isUpdating, setIsUpdating] = useState(false);

    // Delete dialog state
    const [deletingKey, setDeletingKey] = useState<APIKey | null>(null);
    const [deletingKeyName, setDeletingKeyName] = useState("");
    const [isDeleting, setIsDeleting] = useState(false);

    useEffect(() => {
        loadAPIKeys();
    }, []);

    const loadAPIKeys = async () => {
        try {
            setLoading(true);
            setError(null);
            const response = await apiClient.apiKey.listAPIKeys({});
            setApiKeys(response.apiKeys);
        } catch (err) {
            console.error("Failed to load API keys:", err);
            setError("Failed to load API keys. Please try again.");
        } finally {
            setLoading(false);
        }
    };

    const handleCreate = async () => {
        if (!createName.trim()) {
            return;
        }

        try {
            setIsCreating(true);
            setError(null);
            const response = await apiClient.apiKey.createAPIKey({
                name: createName.trim(),
            });

            // Show the created key
            setCreatedKey(response.key);
            setCreateName("");

            // Reload the list
            await loadAPIKeys();
        } catch (err) {
            console.error("Failed to create API key:", err);
            setError("Failed to create API key. Please try again.");
        } finally {
            setIsCreating(false);
        }
    };

    const handleOpenCreateDialog = () => {
        setCreateName("");
        setCreatedKey(null);
        setIsCreateDialogOpen(true);
    };

    const handleCloseCreateDialog = () => {
        setIsCreateDialogOpen(false);
    };

    const handleOpenEditDialog = (key: APIKey) => {
        setEditName(key.name);
        setEditingKey(key);
    };

    const handleCloseEditDialog = () => {
        setEditingKey(null);
    };

    const handleUpdate = async () => {
        if (!editingKey || !editName.trim()) {
            return;
        }

        try {
            setIsUpdating(true);
            setError(null);
            await apiClient.apiKey.updateAPIKey({
                id: editingKey.id,
                name: editName.trim(),
            });

            setEditingKey(null);

            // Reload the list
            await loadAPIKeys();
        } catch (err) {
            console.error("Failed to update API key:", err);
            setError("Failed to update API key. Please try again.");
        } finally {
            setIsUpdating(false);
        }
    };

    const handleOpenDeleteDialog = (key: APIKey) => {
        setDeletingKeyName(key.name);
        setDeletingKey(key);
    };

    const handleCloseDeleteDialog = () => {
        setDeletingKey(null);
    };

    const handleDelete = async () => {
        if (!deletingKey) {
            return;
        }

        try {
            setIsDeleting(true);
            setError(null);
            await apiClient.apiKey.deleteAPIKey({
                id: deletingKey.id,
            });

            // Reload the list
            await loadAPIKeys();

            // Close dialog after successful delete
            setDeletingKey(null);
            setIsDeleting(false);
        } catch (err) {
            console.error("Failed to delete API key:", err);
            setError("Failed to delete API key. Please try again.");
            setIsDeleting(false);
        }
    };

    const formatDate = (timestamp: { seconds: bigint } | undefined) => {
        if (!timestamp) return "Never";
        const date = new Date(Number(timestamp.seconds) * 1000);
        return new Intl.DateTimeFormat("en-US", {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
        }).format(date);
    };

    return (
        <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between mb-8">
                <div>
                    <Heading>API Keys</Heading>
                    <Text className="mt-2">
                        Manage your API keys for programmatic access to your account.
                    </Text>
                </div>
                <Button className="self-start sm:self-auto" color="dark/zinc" onClick={handleOpenCreateDialog}>
                    <PlusIcon/>
                    Create API Key
                </Button>
            </div>

            {error && (
                <div className="mb-6 rounded-lg bg-red-50 p-4 text-sm text-red-800 dark:bg-red-950 dark:text-red-200">
                    {error}
                </div>
            )}

            {loading ? (
                <div className="text-center py-12">
                    <Text>Loading API keys...</Text>
                </div>
            ) : apiKeys.length === 0 ? (
                <div className="text-center py-12 border border-zinc-200 dark:border-zinc-800 rounded-lg">
                    <Text>No API keys yet. Create one to get started.</Text>
                </div>
            ) : (
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableHeader>Name</TableHeader>
                            <TableHeader>Key</TableHeader>
                            <TableHeader>Created</TableHeader>
                            <TableHeader>Last Used</TableHeader>
                            <TableHeader></TableHeader>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {apiKeys.map((key) => (
                            <TableRow key={key.id}>
                                <TableCell className="font-medium">{key.name}</TableCell>
                                <TableCell>
                                    <code className="text-sm text-zinc-600 dark:text-zinc-400">
                                        {key.keyMasked}
                                    </code>
                                </TableCell>
                                <TableCell>{formatDate(key.createdAt)}</TableCell>
                                <TableCell>
                                    {key.lastUsedAt ? (
                                        formatDate(key.lastUsedAt)
                                    ) : (
                                        <Badge color="zinc">Never</Badge>
                                    )}
                                </TableCell>
                                <TableCell>
                                    <div className="flex gap-2 justify-end">
                                        <Button
                                            plain
                                            onClick={() => handleOpenEditDialog(key)}
                                        >
                                            <PencilIcon/>
                                        </Button>
                                        <Button
                                            plain
                                            onClick={() => handleOpenDeleteDialog(key)}
                                        >
                                            <TrashIcon/>
                                        </Button>
                                    </div>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            )}

            {/* Create API Key Dialog */}
            <Dialog open={isCreateDialogOpen} onClose={handleCloseCreateDialog}>
                <DialogTitle>
                    {createdKey ? "API Key Created" : "Create API Key"}
                </DialogTitle>
                <DialogDescription>
                    {createdKey
                        ? "Copy your API key now. You won't be able to see it again."
                        : "Give your API key a name to help you remember what it's for."}
                </DialogDescription>
                <DialogBody>
                    {createdKey ? (
                        <Field>
                            <Label>Your API Key</Label>
                            <Input
                                value={createdKey}
                                readOnly
                                className="font-mono text-sm text-center"
                            />
                        </Field>
                    ) : (
                        <Field>
                            <Label>Name</Label>
                            <Input
                                value={createName}
                                onChange={(e) => setCreateName(e.target.value)}
                                placeholder="e.g., Production API Key"
                                autoFocus
                            />
                        </Field>
                    )}
                </DialogBody>
                <DialogActions>
                    {createdKey ? (
                        <Button onClick={handleCloseCreateDialog}>Done</Button>
                    ) : (
                        <>
                            <Button plain onClick={handleCloseCreateDialog}>
                                Cancel
                            </Button>
                            <Button
                                onClick={handleCreate}
                                disabled={!createName.trim() || isCreating}
                            >
                                {isCreating ? "Creating..." : "Create"}
                            </Button>
                        </>
                    )}
                </DialogActions>
            </Dialog>

            {/* Edit API Key Dialog */}
            <Dialog
                open={editingKey !== null}
                onClose={handleCloseEditDialog}
            >
                <DialogTitle>Rename API Key</DialogTitle>
                <DialogDescription>
                    Update the name of your API key.
                </DialogDescription>
                <DialogBody>
                    <Field>
                        <Label>Name</Label>
                        <Input
                            value={editName}
                            onChange={(e) => setEditName(e.target.value)}
                            placeholder="e.g., Production API Key"
                            autoFocus
                        />
                    </Field>
                </DialogBody>
                <DialogActions>
                    <Button
                        plain
                        onClick={handleCloseEditDialog}
                    >
                        Cancel
                    </Button>
                    <Button
                        onClick={handleUpdate}
                        disabled={!editName.trim() || isUpdating}
                    >
                        {isUpdating ? "Updating..." : "Update"}
                    </Button>
                </DialogActions>
            </Dialog>

            {/* Delete API Key Dialog */}
            <Dialog
                open={deletingKey !== null}
                onClose={handleCloseDeleteDialog}
            >
                <DialogTitle>Delete API Key</DialogTitle>
                <DialogDescription>
                    Are you sure you want to delete "{deletingKeyName}"? This action cannot be undone.
                </DialogDescription>
                <DialogActions>
                    <Button plain onClick={handleCloseDeleteDialog}>
                        Cancel
                    </Button>
                    <Button color="red" onClick={handleDelete} disabled={isDeleting}>
                        {isDeleting ? "Deleting..." : "Delete"}
                    </Button>
                </DialogActions>
            </Dialog>
        </div>
    );
}
