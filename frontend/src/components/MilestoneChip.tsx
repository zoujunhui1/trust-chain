import type { MilestoneState } from '../lib/api'

const LABEL: Record<MilestoneState, string> = {
  0: 'Locked',
  1: 'Released',
  2: 'Proven',
}

const CLASSES: Record<MilestoneState, string> = {
  0: 'bg-locked-tint text-locked',
  1: 'bg-released-tint text-released',
  2: 'bg-proven-tint text-proven',
}

export default function MilestoneChip({ state }: { state: MilestoneState }) {
  return (
    <span className={`rounded-full px-2.5 py-1 text-xs font-medium ${CLASSES[state]}`}>
      {LABEL[state]}
    </span>
  )
}
