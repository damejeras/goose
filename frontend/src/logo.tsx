import iconSvg from '@/assets/icon.svg'

export function Logo({ className, ...props }: React.ComponentPropsWithoutRef<'img'>) {
  return (
    <img
      src={iconSvg}
      alt="Logo"
      data-slot="icon"
      className={className}
      {...props}
    />
  )
}
