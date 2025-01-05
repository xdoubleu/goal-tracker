package helper

import (
	"slices"

	"goal-tracker/api/internal/models"
)

// note: todoist will only use 4 indent levels
// (0: parent, 1: sub, 2: 2*sub, 3: 3*sub, 4: 4*sub).
type GoalTree struct {
	rootGoal *models.Goal
	// contains IDs of all (grand*)children
	childrenIDs []string
	// each key is a subgoal of rootgoal
	subtrees map[string]GoalTree
}

func NewGoalTree() GoalTree {
	return GoalTree{
		rootGoal:    nil,
		childrenIDs: []string{},
		subtrees:    map[string]GoalTree{},
	}
}

func (tree *GoalTree) TryAdd(goal models.Goal) bool {
	if goal.ParentID == nil {
		tree.addNewDirectChild(goal)
		return true
	}

	return tree.walkOverTreesAndTryAdd(goal)
}

func (tree *GoalTree) addNewDirectChild(goal models.Goal) {
	tree.subtrees[goal.ID] = GoalTree{
		rootGoal:    &goal,
		childrenIDs: []string{},
		subtrees:    map[string]GoalTree{},
	}
	tree.addNewIndirectChild(goal)
}

func (tree *GoalTree) addNewIndirectChild(goal models.Goal) {
	tree.childrenIDs = append(tree.childrenIDs, goal.ID)
}

func (tree *GoalTree) walkOverTreesAndTryAdd(goal models.Goal) bool {
	if !tree.hasParent(goal) {
		return false
	}

	for _, subtree := range tree.subtrees {
		// subtree either is parent...
		if subtree.rootGoal.ID == *goal.ParentID {
			// found parent
			subtree.addNewDirectChild(goal)
			return true
		}

		// ... or subtree has parent
		if !subtree.hasParent(goal) {
			continue
		}

		// check subtrees of current tree
		if !subtree.walkOverTreesAndTryAdd(goal) {
			// we shouldn't end up here, the parent should be in one of the trees
			return false
		}

		// this subtree was on the path to the parent
		subtree.addNewIndirectChild(goal)

		return true
	}

	// we shouldn't end up here
	return false
}

func (tree GoalTree) hasParent(goal models.Goal) bool {
	for _, id := range tree.childrenIDs {
		if id == *goal.ParentID {
			return true
		}
	}
	return false
}

type GoalWithSubGoals struct {
	Goal        models.Goal
	HasSubgoals bool
	Subgoals    []GoalWithSubGoals
}

func (tree GoalTree) ToSlice() []GoalWithSubGoals {
	result := []GoalWithSubGoals{}

	for _, subtree := range tree.subtrees {
		goalWithSubGoals := GoalWithSubGoals{
			Goal:        *subtree.rootGoal,
			HasSubgoals: len(subtree.subtrees) > 0,
			Subgoals:    subtree.ToSlice(),
		}

		result = append(result, goalWithSubGoals)
	}

	slices.SortFunc(result, func(a GoalWithSubGoals, b GoalWithSubGoals) int {
		return a.Goal.Order - b.Goal.Order
	})

	return result
}
