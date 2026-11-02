package com.advocatediary.tests;

import org.testng.Assert;
import org.testng.annotations.Test;

public class LasyaAccountFlowTest extends BaseTest {

  @Test
  public void openSignInThenSignInShowsCourtDiary() {
    Assert.assertTrue(signInToDiary().isOpen());
  }
}
