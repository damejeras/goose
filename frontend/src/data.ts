// Mock data for the dashboard

export interface Event {
  id: string
  name: string
  url: string
  date: string
  thumbUrl: string
}

export interface Order {
  id: number
  url: string
  date: string
  amount: {
    usd: string
  }
  customer: {
    name: string
    email: string
  }
  event: Event
}

export async function getEvents(): Promise<Event[]> {
  return [
    {
      id: '1',
      name: 'Bear Hug: Live in Concert',
      url: '/events/1',
      date: 'May 25, 2024',
      thumbUrl: 'https://api.dicebear.com/7.x/shapes/svg?seed=event1',
    },
    {
      id: '2',
      name: 'Viking People',
      url: '/events/2',
      date: 'June 4, 2024',
      thumbUrl: 'https://api.dicebear.com/7.x/shapes/svg?seed=event2',
    },
    {
      id: '3',
      name: 'Six Fingers',
      url: '/events/3',
      date: 'June 12, 2024',
      thumbUrl: 'https://api.dicebear.com/7.x/shapes/svg?seed=event3',
    },
    {
      id: '4',
      name: 'We All Look The Same',
      url: '/events/4',
      date: 'June 19, 2024',
      thumbUrl: 'https://api.dicebear.com/7.x/shapes/svg?seed=event4',
    },
  ]
}

export async function getRecentOrders(): Promise<Order[]> {
  const events = await getEvents()

  return [
    {
      id: 3000,
      url: '/orders/3000',
      date: 'May 9, 2024',
      amount: { usd: '$80.00' },
      customer: {
        name: 'Leslie Alexander',
        email: 'leslie.alexander@example.com',
      },
      event: events[0],
    },
    {
      id: 3001,
      url: '/orders/3001',
      date: 'May 5, 2024',
      amount: { usd: '$299.00' },
      customer: {
        name: 'Michael Foster',
        email: 'michael.foster@example.com',
      },
      event: events[1],
    },
    {
      id: 3002,
      url: '/orders/3002',
      date: 'Apr 28, 2024',
      amount: { usd: '$150.00' },
      customer: {
        name: 'Dries Vincent',
        email: 'dries.vincent@example.com',
      },
      event: events[2],
    },
    {
      id: 3003,
      url: '/orders/3003',
      date: 'Apr 23, 2024',
      amount: { usd: '$80.00' },
      customer: {
        name: 'Lindsay Walton',
        email: 'lindsay.walton@example.com',
      },
      event: events[3],
    },
    {
      id: 3004,
      url: '/orders/3004',
      date: 'Apr 18, 2024',
      amount: { usd: '$199.00' },
      customer: {
        name: 'Courtney Henry',
        email: 'courtney.henry@example.com',
      },
      event: events[0],
    },
    {
      id: 3005,
      url: '/orders/3005',
      date: 'Apr 15, 2024',
      amount: { usd: '$125.00' },
      customer: {
        name: 'Tom Cook',
        email: 'tom.cook@example.com',
      },
      event: events[1],
    },
    {
      id: 3006,
      url: '/orders/3006',
      date: 'Apr 12, 2024',
      amount: { usd: '$89.00' },
      customer: {
        name: 'Whitney Francis',
        email: 'whitney.francis@example.com',
      },
      event: events[2],
    },
    {
      id: 3007,
      url: '/orders/3007',
      date: 'Apr 8, 2024',
      amount: { usd: '$350.00' },
      customer: {
        name: 'Leonard Krasner',
        email: 'leonard.krasner@example.com',
      },
      event: events[3],
    },
  ]
}
