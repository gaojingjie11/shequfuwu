export const GREEN_POINTS_PER_YUAN = 10
export const CENTS_PER_GREEN_POINT = Math.round(100 / GREEN_POINTS_PER_YUAN)

export function getMixedPaymentPreview(amount, greenPoints) {
  const amountInCents = Math.max(0, Math.round(Number(amount || 0) * 100))
  const points = Math.max(0, Math.floor(Number(greenPoints || 0)))
  const maxDeductiblePoints = Math.floor(amountInCents / CENTS_PER_GREEN_POINT)
  const usedPoints = Math.min(points, maxDeductiblePoints)
  const balanceInCents = amountInCents - usedPoints * CENTS_PER_GREEN_POINT

  return {
    points: usedPoints,
    balance: balanceInCents / 100
  }
}
