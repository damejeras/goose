import type {ReactNode} from "react";
import {useLocation, useNavigate} from "react-router-dom";
import {useAuth} from "@/contexts/AuthContext";
import {Avatar} from "@/components/avatar";
import {
    Dropdown,
    DropdownButton,
    DropdownDivider,
    DropdownItem,
    DropdownLabel,
    DropdownMenu,
} from "@/components/dropdown";
import {Navbar, NavbarItem, NavbarSection, NavbarSpacer,} from "@/components/navbar";
import {
    Sidebar,
    SidebarBody,
    SidebarFooter,
    SidebarHeader,
    SidebarItem,
    SidebarLabel,
    SidebarSection,
} from "@/components/sidebar";
import {SidebarLayout} from "@/components/sidebar-layout";
import {ArrowRightStartOnRectangleIcon, ChevronUpIcon,} from "@heroicons/react/16/solid";
import {HomeIcon, KeyIcon,} from "@heroicons/react/20/solid";
import {Logo} from "@/logo";

function AccountDropdownMenu({
                                 anchor,
                                 onSignOut,
                             }: {
    anchor: "top start" | "bottom end";
    onSignOut: () => void;
}) {
    return (
        <DropdownMenu className="min-w-64" anchor={anchor}>
            <DropdownItem href="/api-keys">
                <KeyIcon/>
                <DropdownLabel>API Keys</DropdownLabel>
            </DropdownItem>
            <DropdownDivider/>
            <DropdownItem onClick={onSignOut}>
                <ArrowRightStartOnRectangleIcon/>
                <DropdownLabel>Sign Out</DropdownLabel>
            </DropdownItem>
        </DropdownMenu>
    );
}

export function DashboardLayout({children}: { children: ReactNode }) {
    const location = useLocation();
    const navigate = useNavigate();
    const {user, logout} = useAuth();

    const handleSignOut = async () => {
        await logout();
        navigate("/login");
    };

    // Generate initials from name or email
    const getInitials = (name: string | undefined, email: string | undefined): string => {
        if (name && name.trim()) {
            // Use first letters of first and last name
            const parts = name.trim().split(/\s+/);
            if (parts.length >= 2) {
                return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
            }
            return parts[0][0].toUpperCase();
        }
        // Fallback to first letter of email
        return email ? email[0].toUpperCase() : "U";
    };

    const initials = getInitials(user?.name, user?.email);
    const displayName = user?.name || user?.email.split("@")[0] || "User";

    return (
        <SidebarLayout
            navbar={
                <Navbar>
                    <NavbarSpacer/>
                    <NavbarSection>
                        <Dropdown>
                            <DropdownButton as={NavbarItem}>
                                <Avatar initials={initials} square/>
                            </DropdownButton>
                            <AccountDropdownMenu
                                anchor="bottom end"
                                onSignOut={handleSignOut}
                            />
                        </Dropdown>
                    </NavbarSection>
                </Navbar>
            }
            sidebar={
                <Sidebar>
                    <SidebarHeader>
                        <div className="flex items-center gap-3 px-2 py-2.5">
                            <Logo className="size-6 sm:size-5"/>
                            <span className="text-base/6 font-medium text-zinc-950 sm:text-sm/5 dark:text-white">
                Dashboard
              </span>
                        </div>
                    </SidebarHeader>

                    <SidebarBody>
                        <SidebarSection>
                            <SidebarItem href="/" current={location.pathname === "/"}>
                                <HomeIcon/>
                                <SidebarLabel>Home</SidebarLabel>
                            </SidebarItem>
                        </SidebarSection>
                    </SidebarBody>

                    <SidebarFooter className="max-lg:hidden">
                        <Dropdown>
                            <DropdownButton as={SidebarItem}>
                                <span className="flex min-w-0 items-center gap-3">
                                  <Avatar initials={initials} className="size-10" square alt=""/>
                                  <span className="min-w-0">
                                    <span className="block truncate text-sm/5 font-medium text-zinc-950 dark:text-white">
                                      {displayName}
                                    </span>
                                    <span className="block truncate text-xs/5 font-normal text-zinc-500 dark:text-zinc-400">
                                      {user?.email || "user@example.com"}
                                    </span>
                                  </span>
                                </span>
                                <ChevronUpIcon/>
                            </DropdownButton>
                            <AccountDropdownMenu
                                anchor="top start"
                                onSignOut={handleSignOut}
                            />
                        </Dropdown>
                    </SidebarFooter>
                </Sidebar>
            }
        >
            {children}
        </SidebarLayout>
    );
}
