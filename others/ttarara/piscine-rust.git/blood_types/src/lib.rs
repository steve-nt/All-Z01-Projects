use std::{fmt, str::FromStr};
#[derive(PartialEq, Eq, Hash, Clone, Copy)]pub enum Antigen{A,AB,B,O}
#[derive(PartialEq, Eq, Hash, Clone, Copy)]pub enum RhFactor{Positive,Negative}
#[derive(PartialEq, Eq, Hash, Clone, Copy)]pub struct BloodType{pub antigen:Antigen,pub rh_factor:RhFactor}
impl FromStr for BloodType{type Err=();fn from_str(s:&str)->Result<Self,Self::Err>{
let s=s.trim();let b=s.as_bytes();if b.len()<2{return Err(())}
let rh=b[b.len()-1]as char;let ap=&s[..s.len()-1];
let rh_factor=match rh{'+'=>RhFactor::Positive,'-'=>RhFactor::Negative,_=>return Err(())};
let antigen=match ap{"A"=>Antigen::A,"AB"=>Antigen::AB,"B"=>Antigen::B,"O"=>Antigen::O,_=>return Err(())};
Ok(BloodType{antigen,rh_factor})}}
impl fmt::Debug for BloodType{fn fmt(&self,f:&mut fmt::Formatter<'_>)->fmt::Result{
let a=match self.antigen{Antigen::A=>"A",Antigen::AB=>"AB",Antigen::B=>"B",Antigen::O=>"O"};
let r=if matches!(self.rh_factor,RhFactor::Positive){'+'}else{'-'};
write!(f,"{}{}",a,r)}}
impl BloodType{
const TYPES:[BloodType;8]=[BloodType{antigen:Antigen::AB,rh_factor:RhFactor::Positive},BloodType{antigen:Antigen::AB,rh_factor:RhFactor::Negative},BloodType{antigen:Antigen::O,rh_factor:RhFactor::Positive},BloodType{antigen:Antigen::O,rh_factor:RhFactor::Negative},BloodType{antigen:Antigen::A,rh_factor:RhFactor::Positive},BloodType{antigen:Antigen::A,rh_factor:RhFactor::Negative},BloodType{antigen:Antigen::B,rh_factor:RhFactor::Positive},BloodType{antigen:Antigen::B,rh_factor:RhFactor::Negative}];
pub fn can_receive_from(self,other:Self)->bool{
let abo=match self.antigen{Antigen::A=>matches!(other.antigen,Antigen::A|Antigen::O),Antigen::B=>matches!(other.antigen,Antigen::B|Antigen::O),Antigen::AB=>true,Antigen::O=>other.antigen==Antigen::O};
let rh=match self.rh_factor{RhFactor::Positive=>true,RhFactor::Negative=>other.rh_factor==RhFactor::Negative};
abo&&rh}
pub fn donors(self)->Vec<Self>{Self::TYPES.iter().copied().filter(|d|self.can_receive_from(*d)).collect()}
pub fn recipients(self)->Vec<Self>{Self::TYPES.iter().copied().filter(|r|(*r).can_receive_from(self)).collect()}
}
