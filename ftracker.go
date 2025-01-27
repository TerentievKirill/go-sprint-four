package ftracker

import (
    "fmt"
    "math"
)

// Основные константы, необходимые для расчетов.
const (
    lenStep   = 0.65  // средняя длина шага.
    mInKm     = 1000  // количество метров в километре.
    minInH    = 60    // количество минут в часе.
    kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
    cmInM     = 100   // количество сантиметров в метре.
)

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// action int — количество совершенных действий (число шагов при ходьбе и беге, либо гребков при плавании).
func distance(action int) float64 {
    return float64(action) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 — длительность тренировки в часах.
func meanSpeed(action int, duration float64) float64 {
    if duration == 0 {
        return 0
    }
    distance := distance(action)
    return distance / duration
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// trainingType string — вид тренировки(Бег, Ходьба, Плавание).
// duration float64 — длительность тренировки в часах.
func ShowTrainingInfo(action int, trainingType string, duration, weight, height float64, lengthPool, countPool int) string {
    switch trainingType {
    case "Бег":
        distance := distance(action) // получаем дистанцию
        speed := meanSpeed(action, duration) // получаем среднюю скорость
        calories := RunningSpentCalories(action, weight, duration) // получаем количество сожженных калорий
        return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, duration, distance, speed, calories)
    case "Ходьба":
        distance := distance(action) // получаем дистанцию
        speed := meanSpeed(action, duration) // получаем среднюю скорость
        calories := WalkingSpentCalories(action, duration, weight, height) // получаем количество сожженных калорий
        return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, duration, distance, speed, calories)
    case "Плавание":
        distance := distance(action) // получаем дистанцию
        speed := swimmingMeanSpeed(lengthPool, countPool, duration) // получаем среднюю скорость
        calories := SwimmingSpentCalories(lengthPool, countPool, duration, weight) // получаем количество сожженных калорий
        return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainingType, duration, distance, speed, calories)
    default:
        return "неизвестный тип тренировки"
    }
}



// Константы для расчета калорий, расходуемых при беге.
const (
    runningCaloriesMeanSpeedMultiplier = 18   // множитель средней скорости.
    runningCaloriesMeanSpeedShift      = 1.79 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// weight float64 — вес пользователя.
// duration float64 — длительность тренировки в часах.
func RunningSpentCalories(action int, weight, duration float64) float64 {
    const (
        runningCaloriesMeanSpeedMultiplier = 18
        runningCaloriesMeanSpeedShift      = 1.79
    )
    speed := meanSpeed(action, duration)  // получаем среднюю скорость
    calories := ((runningCaloriesMeanSpeedMultiplier * speed * runningCaloriesMeanSpeedShift) * weight / mInKm * duration * minInH)
    return calories
}


// Константы для расчета калорий, расходуемых при ходьбе.
const (
    walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
    walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 — длительность тренировки в часах.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(action int, duration, weight, height float64) float64 {
    const (
        walkingCaloriesWeightMultiplier = 0.035  // множитель массы тела
        walkingSpeedHeightMultiplier    = 0.029  // множитель скорости/роста
    )

    // Формула для расчета калорий при ходьбе с учетом правильного порядка операций
    calories := ((walkingCaloriesWeightMultiplier*weight + (math.Pow(meanSpeed(action, duration)*kmhInMsec, 2)/(height/cmInM))*walkingSpeedHeightMultiplier*weight) * duration * minInH)

    return calories
}


// Константы для расчета калорий, расходуемых при плавании.
const (
    swimmingCaloriesMeanSpeedShift   = 1.1  // среднее количество сжигаемых колорий при плавании относительно скорости.
    swimmingCaloriesWeightMultiplier = 2    // множитель веса при плавании.
)

// swimmingMeanSpeed возвращает среднюю скорость при плавании.
//
// Параметры:
//
// lengthPool int — длина бассейна в метрах.
// countPool int — сколько раз пользователь переплыл бассейн.
// duration float64 — длительность тренировки в часах.
func swimmingMeanSpeed(lengthPool, countPool int, duration float64) float64 {
    if duration == 0 {
        return 0
    }
    return float64(lengthPool) * float64(countPool) / mInKm / duration
}

// SwimmingSpentCalories возвращает количество потраченных калорий при плавании.
//
// Параметры:
//
// lengthPool int — длина бассейна в метрах.
// countPool int — сколько раз пользователь переплыл бассейн.
// duration float64 — длительность тренировки в часах.
// weight float64 — вес пользователя.
func SwimmingSpentCalories(lengthPool, countPool int, duration, weight float64) float64 {
    const (
        swimmingCaloriesMeanSpeedShift   = 1.1
        swimmingCaloriesWeightMultiplier = 2
    )
    speed := swimmingMeanSpeed(lengthPool, countPool, duration)  // получаем среднюю скорость плавания
    calories := (speed + swimmingCaloriesMeanSpeedShift) * swimmingCaloriesWeightMultiplier * weight * duration
    return calories
}

