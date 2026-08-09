export const pageSize = 20

interface FlipResult {
    duration: number;
    steps: number;
}

// 根据判断结果返回true/false，如果判断通过，则在调用函数以后、一直到该作用域结束，obj将被视为`FlipResult`类型
export function isFlipResult(obj: any): obj is FlipResult {
    return (
        typeof obj === "object" &&
        obj !== null &&
        typeof obj.duration === "number" &&
        typeof obj.steps === "number"
    )
}
