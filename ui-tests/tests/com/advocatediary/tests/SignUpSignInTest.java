package com.advocatediary.tests;

import com.advocatediary.pages.SignUpPage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class SignUpSignInTest extends BaseTest {

  @Test
  public void existingLawyerSignsInFromWelcomeBack() {
    Assert.assertTrue(signInToDiary().isOpen());
  }

  @Test
  public void shortPasswordShowsHintOnSignUp() {
    SignUpPage signUp = new SignUpPage(driver).open(baseUrl);
    signUp.submitInvalid("Test Advocate", "short@chamber.test", "short");
    Assert.assertTrue(signUp.showsPasswordHint());
  }
}
