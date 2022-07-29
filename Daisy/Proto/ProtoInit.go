package Proto

import (
	"Cinder/Base/ProtoDef"
)

func init() {
	ProtoDef.AddDef(431, &RspOfflineAwardData{})
	ProtoDef.AddDef(432, &OfflineAwardData{})
	ProtoDef.AddDef(433, &OfflineAwardItem{})
	ProtoDef.AddDef(434, &ChatMessage{})
	ProtoDef.AddDef(435, &ChatHistoryMessage{})
	ProtoDef.AddDef(436, &ChannelChatHistoryMessage{})
	ProtoDef.AddDef(437, &PVector3{})
	ProtoDef.AddDef(438, &PVector3Array{})
	ProtoDef.AddDef(439, &Int32Map{})
	ProtoDef.AddDef(440, &Int32StringMap{})
	ProtoDef.AddDef(441, &Int32Array{})
	ProtoDef.AddDef(442, &StringArry{})
	ProtoDef.AddDef(443, &PawnType{})
	ProtoDef.AddDef(444, &Camp{})
	ProtoDef.AddDef(445, &ActionType{})
	ProtoDef.AddDef(446, &AttrType{})
	ProtoDef.AddDef(447, &StatType{})
	ProtoDef.AddDef(448, &AttackSrc{})
	ProtoDef.AddDef(449, &AttackShapeType{})
	ProtoDef.AddDef(450, &DamageType{})
	ProtoDef.AddDef(451, &MoveMode{})
	ProtoDef.AddDef(452, &HitType{})
	ProtoDef.AddDef(453, &SkillState{})
	ProtoDef.AddDef(454, &SkillBreakReason{})
	ProtoDef.AddDef(455, &SkillEndReason{})
	ProtoDef.AddDef(456, &FightLogType{})
	ProtoDef.AddDef(457, &DamageFloatWordType{})
	ProtoDef.AddDef(458, &FightAwardData{})
	ProtoDef.AddDef(459, &FightResult{})
	ProtoDef.AddDef(460, &DebugSceneInfo{})
	ProtoDef.AddDef(461, &Position{})
	ProtoDef.AddDef(462, &FightRoleInfo{})
	ProtoDef.AddDef(463, &FightAIInfo{})
	ProtoDef.AddDef(464, &FightNpcInfo{})
	ProtoDef.AddDef(465, &PawnInherit{})
	ProtoDef.AddDef(466, &FightInherit{})
	ProtoDef.AddDef(467, &PawnInfo{})
	ProtoDef.AddDef(468, &FightReplay{})
	ProtoDef.AddDef(469, &FightFrame{})
	ProtoDef.AddDef(470, &FightAction{})
	ProtoDef.AddDef(471, &SummonPawn{})
	ProtoDef.AddDef(472, &MoveBegin{})
	ProtoDef.AddDef(473, &MoveEnd{})
	ProtoDef.AddDef(474, &FixMoveData{})
	ProtoDef.AddDef(475, &UseSkill{})
	ProtoDef.AddDef(476, &BreakSkill{})
	ProtoDef.AddDef(477, &NewAttack{})
	ProtoDef.AddDef(478, &DelAttack{})
	ProtoDef.AddDef(479, &AttackHit{})
	ProtoDef.AddDef(480, &BeHit{})
	ProtoDef.AddDef(481, &AddBuff{})
	ProtoDef.AddDef(482, &RemoveBuff{})
	ProtoDef.AddDef(483, &ChangeAttr{})
	ProtoDef.AddDef(484, &ChangeStat{})
	ProtoDef.AddDef(485, &FightBegin{})
	ProtoDef.AddDef(486, &FightEnd{})
	ProtoDef.AddDef(487, &DebugInfo{})
	ProtoDef.AddDef(488, &ChangeSkillState{})
	ProtoDef.AddDef(489, &CombineSkillEndTime{})
	ProtoDef.AddDef(490, &CombineSkillPoint{})
	ProtoDef.AddDef(491, &SetTarget{})
	ProtoDef.AddDef(492, &AttackAoeTrans{})
	ProtoDef.AddDef(493, &AttackAoeShape{})
	ProtoDef.AddDef(494, &AttackShowAoe{})
	ProtoDef.AddDef(495, &AttackMoveAoe{})
	ProtoDef.AddDef(496, &SpawnMonsterInfo{})
	ProtoDef.AddDef(497, &Friend{})
	ProtoDef.AddDef(498, &FriendListArray{})
	ProtoDef.AddDef(499, &FriendArray{})
	ProtoDef.AddDef(500, &FriendListData{})
	ProtoDef.AddDef(501, &ItemBase{})
	ProtoDef.AddDef(502, &ItemEnum{})
	ProtoDef.AddDef(503, &Item{})
	ProtoDef.AddDef(504, &OwnerInfo{})
	ProtoDef.AddDef(505, &Equipment{})
	ProtoDef.AddDef(506, &AffixData{})
	ProtoDef.AddDef(507, &SkillItem{})
	ProtoDef.AddDef(508, &ItemExpand{})
	ProtoDef.AddDef(509, &ItemContainer{})
	ProtoDef.AddDef(510, &ContainerEnum{})
	ProtoDef.AddDef(511, &Items{})
	ProtoDef.AddDef(512, &ItemsMap{})
	ProtoDef.AddDef(513, &DropMaterial{})
	ProtoDef.AddDef(514, &GetItemData{})
	ProtoDef.AddDef(515, &Mail{})
	ProtoDef.AddDef(516, &MailAttachment{})
	ProtoDef.AddDef(517, &MailBox{})
	ProtoDef.AddDef(518, &MailAwardItems{})
	ProtoDef.AddDef(519, &RedPointInfo{})
	ProtoDef.AddDef(520, &RoleCache{})
	ProtoDef.AddDef(521, &Role{})
	ProtoDef.AddDef(522, &RoleBase{})
	ProtoDef.AddDef(523, &RoleInviteInfo{})
	ProtoDef.AddDef(524, &RoleChat{})
	ProtoDef.AddDef(525, &ShareSpoils{})
	ProtoDef.AddDef(526, &GiveSkillCount{})
	ProtoDef.AddDef(527, &RequestSkill{})
	ProtoDef.AddDef(528, &FastBattle{})
	ProtoDef.AddDef(529, &FastBattleStageInfo{})
	ProtoDef.AddDef(530, &SupplyInfo{})
	ProtoDef.AddDef(531, &State{})
	ProtoDef.AddDef(532, &SupplyAwardItem{})
	ProtoDef.AddDef(533, &SupplyAwardData{})
	ProtoDef.AddDef(534, &Title{})
	ProtoDef.AddDef(535, &TitleInfo{})
	ProtoDef.AddDef(536, &TeamSeasonInfo{})
	ProtoDef.AddDef(537, &RoleSeasonInfo{})
	ProtoDef.AddDef(538, &RankData{})
	ProtoDef.AddDef(539, &RankTeamData{})
	ProtoDef.AddDef(540, &RankMemberData{})
	ProtoDef.AddDef(541, &SpecialAgent{})
	ProtoDef.AddDef(542, &SpecialAgentBase{})
	ProtoDef.AddDef(543, &SpecialAgentTalent{})
	ProtoDef.AddDef(544, &BuildData{})
	ProtoDef.AddDef(545, &BuildSkillData{})
	ProtoDef.AddDef(546, &SkillData{})
	ProtoDef.AddDef(547, &FightAttr{})
	ProtoDef.AddDef(548, &TalentData{})
	ProtoDef.AddDef(549, &TeamCache{})
	ProtoDef.AddDef(550, &Team{})
	ProtoDef.AddDef(551, &TeamBase{})
	ProtoDef.AddDef(552, &TeamMemberInfo{})
	ProtoDef.AddDef(553, &TeamRun{})
	ProtoDef.AddDef(554, &RaidInfo{})
	ProtoDef.AddDef(555, &GuideInfo{})
	ProtoDef.AddDef(556, &TeamApplyInfo{})
	ProtoDef.AddDef(557, &TeamPart{})
	ProtoDef.AddDef(558, &RoleArry{})
	ProtoDef.AddDef(559, &Supply{})
	ProtoDef.AddDef(560, &Recruitment{})
	ProtoDef.AddDef(561, &RecruitmentMember{})
	ProtoDef.AddDef(562, &Recruitments{})
	ProtoDef.AddDef(563, &Loner{})
	ProtoDef.AddDef(564, &Loners{})
	ProtoDef.AddDef(565, &Applyer{})
	ProtoDef.AddDef(566, &Applyers{})
	ProtoDef.AddDef(567, &Inviter{})
	ProtoDef.AddDef(568, &Inviters{})
	ProtoDef.InitProtoDefData()

}
