// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"math/rand"
	"sort"
)

func LogSummary(gradeCounts map[string]int, toppedCounts map[string]int, showCongrats bool, lastAttemptDifficulty int) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if showCongrats {
			if lastAttemptDifficulty <= 4 {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h2 class=\"text-3xl font-bold mb-4 text-gray-900 dark:text-gray-100\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var2 string
				templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(getRandomHypeText())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/log_summary.templ`, Line: 12, Col: 93}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h2>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h2 class=\"text-3xl font-bold mb-4 text-gray-900 dark:text-gray-100\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(getRandomPepText())
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/log_summary.templ`, Line: 14, Col: 92}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h2>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		} else if len(gradeCounts) == 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h2 class=\"text-3xl font-bold mb-4 text-gray-900 dark:text-gray-100\">Let's Get Started!</h2><p class=\"text-xl mb-6 text-gray-700 dark:text-gray-300\">Ready to crush some boulders? Log your first attempt of the day!</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h3 class=\"text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100\">Today's Progress</h3><div class=\"grid grid-cols-3 gap-2 mb-6\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(gradeCounts) > 0 {
			for _, grade := range sortedGrades(gradeCounts) {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bg-gray-100 dark:bg-gray-900 p-4 rounded flex items-center\"><span class=\"font-bold text-gray-800 dark:text-gray-200\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(grade)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/log_summary.templ`, Line: 25, Col: 69}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(":</span> <span class=\"text-green-600 dark:text-green-400\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var5 string
				templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(templ.EscapeString(fmt.Sprintf("%d", toppedCounts[grade])))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/log_summary.templ`, Line: 25, Col: 189}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span>/<span class=\"text-gray-800 dark:text-gray-200\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(templ.EscapeString(fmt.Sprintf("%d", gradeCounts[grade])))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/log_summary.templ`, Line: 25, Col: 305}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"col-span-2 text-gray-500 dark:text-gray-400\">No boulders logged today</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><button hx-get=\"/log/grade\" hx-target=\"#main-content\" hx-swap=\"innerHTML\" hx-push-url=\"true\" class=\"bg-blue-500 hover:bg-blue-600 dark:bg-blue-600 dark:hover:bg-blue-700 text-white font-bold py-3 px-6 rounded\">Log New Boulder</button>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func getRandomHypeText() string {
	hypeTexts := []string{
		"Congratulations! You're crushing it!",
		"Amazing work! Keep pushing those limits!",
		"You're on fire today! Great job!",
		"Incredible progress! You're a rock star!",
		"Awesome effort! You're scaling new heights!",
		"Fantastic climbing! You're reaching new peaks!",
		"Superb performance! You're defying gravity!",
		"Outstanding work! You're conquering those boulders!",
		"Phenomenal climbing! You're unstoppable!",
		"Brilliant effort! You're climbing to new heights!",
		"You're a rock star! Crushing boulders like a pro!",
		"Gravity called, it wants its rules back! You're defying physics!",
		"Holy chalk bucket! You're on fire today!",
		"Sending it like there's no tomorrow! You're unstoppable!",
		"You've got the touch! That boulder didn't stand a chance!",
		"Boom! Another one bites the dust! You're a boulder-slaying machine!",
		"Chalk up that victory! You're climbing like a legend!",
		"Who needs stairs when you can climb walls? You're crushing it!",
		"You're making these boulders look like pebbles! Incredible!",
		"Spiderman called, he wants climbing tips from you!",
		"You're not just sending it, you're express shipping it! Amazing!",
		"Forget the elevator, you're taking the rock face express!",
		"You're not climbing, you're flying up these boulders!",
		"The mountain goats are jealous of your skills! Keep it up!",
		"You're leaving no stone unturned... or unclimbed! Fantastic!",
		"Chalk-covered hands, victory-covered climbs! You're on a roll!",
		"You're not just reaching new heights, you're redefining them!",
		"Boulders beware! You're the new sheriff in town!",
		"You're writing your name in the stars... of the climbing world!",
		"Move over, gravity! There's a new force to be reckoned with!",
		"You're not just climbing, you're rewriting the laws of physics!",
		"Forget walking on sunshine, you're climbing on stardust!",
		"You're turning these boulders into your personal playground!",
		"The rock whisperer strikes again! You're in perfect harmony with the wall!",
		"You're not breaking records, you're shattering them! Incredible performance!",
	}
	return hypeTexts[rand.Intn(len(hypeTexts))]
}

func getRandomPepText() string {
	pepTexts := []string{
		"Great attempt! Every try makes you stronger!",
		"So close! You're making progress with each climb!",
		"Keep pushing! You're building strength and technique!",
		"Awesome effort! The send is just around the corner!",
		"You're getting there! Each attempt brings you closer to success!",
		"Stay motivated! Your persistence will pay off!",
		"Fantastic try! You're learning with every move!",
		"Don't give up! You're developing crucial skills!",
		"Impressive effort! Your determination is inspiring!",
		"Keep at it! Success comes from embracing the challenge!",
		"Rome wasn't built in a day, and neither is climbing prowess. Keep at it!",
		"Every attempt is a step towards success. You're making progress!",
		"The only way is up, and you're on your way! Keep pushing!",
		"Boulders are just puzzles waiting to be solved. You've got this!",
		"Remember, even the pros fall sometimes. You're doing great!",
		"Your determination is as solid as the rock you're climbing. Keep going!",
		"Each attempt makes you stronger. You're building climbing superpowers!",
		"The boulder may have won this round, but the war isn't over. You've got this!",
		"Chalk up and try again! Success is just another attempt away!",
		"You're not falling, you're just practicing your aerial moves. Keep at it!",
		"The boulder is testing you, and you're rising to the challenge!",
		"Every fall is a lesson in disguise. You're getting smarter with each try!",
		"You're not stuck, you're just gathering data for your next attempt!",
		"The send train is coming, and you've got a first-class ticket!",
		"You're not just climbing, you're having a conversation with the rock. Keep talking!",
		"Gravity is just a suggestion, and you're learning to ignore it. Keep pushing!",
		"You're turning sweat into progress. Each attempt is making you stronger!",
		"The boulder is tough, but you're tougher. Show it who's boss!",
		"You're not failing, you're finding ways that don't work yet. Keep exploring!",
		"Every almost-send is a future victory in the making. You're so close!",
		"You're not just climbing, you're dancing with the rock. Keep grooving!",
		"The boulder might be stubborn, but you're more persistent. Don't give up!",
		"You're collecting beta with every attempt. Soon, you'll crack the code!",
		"Remember, diamonds are made under pressure. You're becoming a gem of a climber!",
		"You're not just trying, you're redefining your limits. Keep pushing those boundaries!",
	}
	return pepTexts[rand.Intn(len(pepTexts))]

}

func sortedGrades(gradeCounts map[string]int) []string {
	grades := make([]string, 0, len(gradeCounts))
	for grade := range gradeCounts {
		grades = append(grades, grade)
	}
	sort.Strings(grades)
	return grades
}

var _ = templruntime.GeneratedTemplate
