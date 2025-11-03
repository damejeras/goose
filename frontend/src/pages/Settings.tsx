import { Button } from '@/components/button'
import { Checkbox, CheckboxField } from '@/components/checkbox'
import { Divider } from '@/components/divider'
import { Heading, Subheading } from '@/components/heading'
import { Input } from '@/components/input'
import { Select } from '@/components/select'
import { Text } from '@/components/text'
import { Textarea } from '@/components/textarea'
import { Label } from '@/components/fieldset'
import { Switch, SwitchField } from '@/components/switch'

export default function Settings() {
  return (
    <form method="post" className="mx-auto max-w-4xl">
      <Heading>Settings</Heading>
      <Divider className="my-10 mt-6" />

      <section className="grid gap-x-8 gap-y-6 sm:grid-cols-2">
        <div className="space-y-1">
          <Subheading>Organization Name</Subheading>
          <Text>This will be displayed on your public profile.</Text>
        </div>
        <div>
          <Input aria-label="Organization Name" name="name" defaultValue="SaaS Dashboard" />
        </div>
      </section>

      <Divider className="my-10" soft />

      <section className="grid gap-x-8 gap-y-6 sm:grid-cols-2">
        <div className="space-y-1">
          <Subheading>Organization Bio</Subheading>
          <Text>This will be displayed on your public profile. Maximum 240 characters.</Text>
        </div>
        <div>
          <Textarea
            aria-label="Organization Bio"
            name="bio"
            rows={3}
            defaultValue="Building the next generation of SaaS products with modern tools and frameworks."
          />
        </div>
      </section>

      <Divider className="my-10" soft />

      <section className="grid gap-x-8 gap-y-6 sm:grid-cols-2">
        <div className="space-y-1">
          <Subheading>Contact Email</Subheading>
          <Text>This is how customers can contact you for support.</Text>
        </div>
        <div className="space-y-4">
          <Input
            type="email"
            aria-label="Contact Email"
            name="email"
            defaultValue="support@saas-dashboard.com"
          />
          <CheckboxField>
            <Checkbox name="email_is_public" defaultChecked />
            <Label>Show email on public profile</Label>
          </CheckboxField>
        </div>
      </section>

      <Divider className="my-10" soft />

      <section className="grid gap-x-8 gap-y-6 sm:grid-cols-2">
        <div className="space-y-1">
          <Subheading>Language & Region</Subheading>
          <Text>Set your preferred language and timezone.</Text>
        </div>
        <div className="space-y-4">
          <div>
            <Text className="font-medium">Language</Text>
            <Select aria-label="Language" name="language" defaultValue="en" className="mt-2">
              <option value="en">English</option>
              <option value="es">Spanish</option>
              <option value="fr">French</option>
              <option value="de">German</option>
            </Select>
          </div>
          <div>
            <Text className="font-medium">Timezone</Text>
            <Select aria-label="Timezone" name="timezone" defaultValue="pst" className="mt-2">
              <option value="pst">Pacific Standard Time</option>
              <option value="est">Eastern Standard Time</option>
              <option value="cst">Central Standard Time</option>
              <option value="mst">Mountain Standard Time</option>
            </Select>
          </div>
        </div>
      </section>

      <Divider className="my-10" soft />

      <section className="grid gap-x-8 gap-y-6 sm:grid-cols-2">
        <div className="space-y-1">
          <Subheading>Notifications</Subheading>
          <Text>Manage how you receive notifications.</Text>
        </div>
        <div className="space-y-6">
          <SwitchField>
            <Label>Email notifications</Label>
            <Text>Receive email updates about your account activity.</Text>
            <Switch name="email_notifications" defaultChecked />
          </SwitchField>
          <SwitchField>
            <Label>Push notifications</Label>
            <Text>Receive push notifications on your devices.</Text>
            <Switch name="push_notifications" />
          </SwitchField>
          <SwitchField>
            <Label>Weekly digest</Label>
            <Text>Receive a weekly summary of your account activity.</Text>
            <Switch name="weekly_digest" defaultChecked />
          </SwitchField>
        </div>
      </section>

      <Divider className="my-10" soft />

      <section className="grid gap-x-8 gap-y-6 sm:grid-cols-2">
        <div className="space-y-1">
          <Subheading>Privacy</Subheading>
          <Text>Control your privacy and data settings.</Text>
        </div>
        <div className="space-y-4">
          <CheckboxField>
            <Checkbox name="analytics" defaultChecked />
            <Label>Allow analytics tracking</Label>
          </CheckboxField>
          <CheckboxField>
            <Checkbox name="marketing" />
            <Label>Receive marketing emails</Label>
          </CheckboxField>
          <CheckboxField>
            <Checkbox name="third_party" />
            <Label>Share data with third-party integrations</Label>
          </CheckboxField>
        </div>
      </section>

      <Divider className="my-10" soft />

      <div className="flex justify-end gap-4">
        <Button type="reset" plain>
          Reset
        </Button>
        <Button type="submit">Save changes</Button>
      </div>
    </form>
  )
}
