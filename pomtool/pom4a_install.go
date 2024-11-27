package pomtool

import (
	"fmt"
	"strings"
)

type pom4aInstall struct {
	pom4aBase
}

func (inst *pom4aInstall) run() error {

	steps := make([]func() error, 0)

	steps = append(steps, inst.stepLoadWD)
	steps = append(steps, inst.stepFindPOM)
	steps = append(steps, inst.stepLoadPOM)
	steps = append(steps, inst.stepCheckPOM)
	steps = append(steps, inst.stepLocateRepo)
	steps = append(steps, inst.stepFindAAR1)
	steps = append(steps, inst.stepLocateAAR2)

	steps = append(steps, inst.stepCopyFiles)
	steps = append(steps, inst.stepDone)

	for _, step := range steps {
		err := step()
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *pom4aInstall) computePom2Path() File {
	tmp := inst.aar2
	name1 := tmp.Name()
	dir := tmp.Parent()
	idx := strings.LastIndexByte(name1, '.')
	name2 := name1[0:idx] + ".pom"
	return dir.Child(name2)
}

func (inst *pom4aInstall) createCopyTask(src, dst File) *copyFileTask {
	return &copyFileTask{src: src, dst: dst}
}

func (inst *pom4aInstall) stepCopyFiles() error {

	pom2 := inst.computePom2Path()
	all := make([]*copyFileTask, 0)

	all = append(all, inst.createCopyTask(inst.aar1, inst.aar2))
	all = append(all, inst.createCopyTask(inst.pom, pom2))

	for _, t := range all {
		err := t.apply()
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *pom4aInstall) stepDone() error {

	fmt.Println("group    : ", inst.groupID)
	fmt.Println("artifact : ", inst.artifactID)
	fmt.Println("version  : ", inst.version)
	fmt.Println()

	fmt.Println("pom        : ", inst.pom)
	fmt.Println("repository : ", inst.repository)
	fmt.Println("aar2       : ", inst.aar2)
	fmt.Println("aar1       : ", inst.aar1)

	fmt.Println()
	fmt.Println("Done.")
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type copyFileTask struct {
	src File
	dst File
}

func (inst *copyFileTask) apply() error {

	src := inst.src
	dst := inst.dst
	if dst.Exists() {
		fmt.Println("[warn] the dest file exist, skip to copy it. file = ", dst)
		return nil
	}

	dir := dst.Parent()
	if !dir.Exists() {
		dir.Mkdirs()
	}

	return src.CopyTo(dst)
}
