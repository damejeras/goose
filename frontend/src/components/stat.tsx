import { Divider } from '@/components/divider'
import clsx from 'clsx'

export function Stat({
  title,
  value,
  change,
}: {
  title: string
  value: string
  change: string
}) {
  const isPositive = change.startsWith('+')

  return (
    <div>
      <Divider />
      <div className="mt-6 text-lg/6 font-medium sm:text-sm/6">{title}</div>
      <div className="mt-3 text-3xl/8 font-semibold sm:text-2xl/8">{value}</div>
      <div className="mt-3 text-sm/6 sm:text-xs/6">
        <span
          className={clsx(
            'inline-flex items-center gap-x-1.5 rounded-md px-1.5 py-0.5 text-sm/5 font-medium sm:text-xs/5 forced-colors:outline',
            isPositive
              ? 'bg-lime-400/20 text-lime-700 dark:bg-lime-400/10 dark:text-lime-300'
              : 'bg-pink-400/15 text-pink-700 dark:bg-pink-400/10 dark:text-pink-400'
          )}
        >
          {change}
        </span>
        <span className="text-zinc-500"> from last week</span>
      </div>
    </div>
  )
}
