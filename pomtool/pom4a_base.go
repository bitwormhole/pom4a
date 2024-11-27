package pomtool

import (
	"fmt"
	"os"
	"os/user"
	"strings"
)

type pom4aBase struct {
	wd         File
	pom        File // file:pom.xml
	aar1       File // AAR source
	aar2       File // AAR in repo
	m2         File // dir:~/.m2
	repository File // dir:~/.m2/repository

	properties PomProperties

	artifactID   string
	groupID      string
	version      string
	modelVersion string // must be '4.0.0'
	packaging    string // must be 'aar'
}

func (inst *pom4aBase) stepLoadWD() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	inst.wd = File(wd)
	return nil
}

func (inst *pom4aBase) stepFindAAR1() error {
	path2 := inst.properties.GetProperty("apt.outputs.aar")
	inst.aar1 = inst.pom.Parent().Child(path2)
	return nil
}

func (inst *pom4aBase) stepLocateAAR2() error {

	builder := &strings.Builder{}
	builder.WriteString(inst.repository.String())

	name := inst.artifactID + "-" + inst.version + ".aar"
	elist := strings.Split(inst.groupID, ".")

	elist = append(elist, inst.artifactID)
	elist = append(elist, inst.version)
	elist = append(elist, name)

	for _, item := range elist {
		builder.WriteRune('/')
		builder.WriteString(item)
	}
	a2 := File(builder.String())
	inst.aar2 = a2.Normalize()
	return nil
}

func (inst *pom4aBase) stepLocateRepo() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	home := File(u.HomeDir)
	m2 := home.Child(".m2")
	repo := m2.Child("repository")
	inst.m2 = m2
	inst.repository = repo
	return nil
}

func (inst *pom4aBase) stepFindPOM() error {
	wd := inst.wd
	inst.pom = wd.Child("pom.xml")
	return nil
}

func (inst *pom4aBase) stepLoadPOM() error {
	data, err := inst.pom.ReadBinary()
	if err != nil {
		return err
	}
	pom := &POM{}
	project, err := pom.Parse(data)
	if err != nil {
		return err
	}
	inst.groupID = project.GroupID
	inst.artifactID = project.ArtifactID
	inst.version = project.Version
	inst.properties = project.Properties
	inst.modelVersion = project.ModelVersion
	inst.packaging = project.Packaging
	return nil
}

func (inst *pom4aBase) stepCheckPOM() error {
	const (
		wantMV   = "4.0.0"
		wantPack = "aar"
	)
	mv := inst.modelVersion
	pack := inst.packaging

	if mv != wantMV {
		return fmt.Errorf("pom.modelVersion MUST be '%s'", wantMV)
	}

	if pack != wantPack {
		return fmt.Errorf("pom.packaging MUST be '%s'", wantPack)
	}

	return nil
}
