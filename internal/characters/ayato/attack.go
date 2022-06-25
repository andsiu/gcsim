package ayato

import (
	"fmt"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var attackFrames [][]int
var attackHitmarks = [][]int{{12}, {18}, {20}, {22, 25}, {41}}
var shunsuikenFrames []int

const normalHitNum = 5
const shunsuikenHitmark = 5

func init() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 24)
	attackFrames[0][action.ActionAttack] = 15

	// TODO: charge cancels are missing?
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 27)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 30)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 27)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 63)

	// NA (in skill) -> x
	shunsuikenFrames = frames.InitNormalCancelSlice(shunsuikenHitmark, 23)
}

func (c *char) Attack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		Abil:       fmt.Sprintf("Normal %v", c.NormalCounter),
		ActorIndex: c.Index,
		AttackTag:  combat.AttackTagNormal,
		ICDTag:     combat.ICDTagNormalAttack,
		ICDGroup:   combat.ICDGroupDefault,
		Element:    attributes.Physical,
		Durability: 25,
	}

	if c.Core.Status.Duration("soukaikanka") > 0 {
		ai.Mult = shunsuiken[c.NormalCounter][c.TalentLvlSkill()]
		c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, shunsuikenHitmark, c.generateParticles, c.skillStacks)
	} else {
		for i, mult := range attack[c.NormalCounter] {
			ai.Mult = mult[c.TalentLvlAttack()]
			c.Core.QueueAttack(
				ai,
				combat.NewDefSingleTarget(1, combat.TargettableEnemy),
				attackHitmarks[c.NormalCounter][i],
				attackHitmarks[c.NormalCounter][i],
			)
		}
	}

	defer c.AdvanceNormalIndex()

	// while in skill
	if c.Core.Status.Duration("soukaikanka") > 0 {
		return action.ActionInfo{
			Frames: func(next action.Action) int {
				return frames.AtkSpdAdjust(shunsuikenFrames[next], c.Stat(attributes.AtkSpd))
			},
			AnimationLength: shunsuikenFrames[action.InvalidAction],
			CanQueueAfter:   shunsuikenHitmark,
			State:           action.NormalAttackState,
		}
	}

	// normal state
	return action.ActionInfo{
		Frames:          frames.NewAttackFunc(c.Character, attackFrames),
		AnimationLength: attackFrames[c.NormalCounter][action.InvalidAction],
		CanQueueAfter:   attackHitmarks[c.NormalCounter][len(attackHitmarks[c.NormalCounter])-1],
		State:           action.NormalAttackState,
	}
}